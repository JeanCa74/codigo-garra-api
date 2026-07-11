'use strict';

/* ════════════════════════════════════════════════════════════════════════
   Código Garra — interfaz de demostración
   Seguridad aplicada en el cliente:
   - El JWT vive en sessionStorage (se borra al cerrar la pestaña).
   - Todo dato de la API se pinta con textContent → imposible inyectar HTML (XSS).
   - Si el token expira o es inválido (401) se cierra la sesión automáticamente.
   - Los botones de eliminar solo se muestran al rol admin; el servidor
     igualmente valida el rol (403), la UI solo acompaña.
   ════════════════════════════════════════════════════════════════════════ */

const API = '/api/v1';
const CLAVE_TOKEN = 'garra_token';

let token = null;
let usuario = null; // { email, rol, exp }
let pestanaActiva = 'panel';

/* ══════════ Helpers de DOM ══════════ */

const $ = (sel) => document.querySelector(sel);

// el('td', 'clase', 'texto') — crea elementos usando SIEMPRE textContent.
function el(etiqueta, clase, texto) {
  const nodo = document.createElement(etiqueta);
  if (clase) nodo.className = clase;
  if (texto !== undefined && texto !== null) nodo.textContent = String(texto);
  return nodo;
}

function limpiar(nodo) {
  while (nodo.firstChild) nodo.removeChild(nodo.firstChild);
}

let toastTimer = null;
function toast(mensaje, tipo) {
  const t = $('#toast');
  t.textContent = mensaje;
  t.className = 'toast' + (tipo ? ' ' + tipo : '');
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => t.classList.add('oculto'), 3500);
}

/* ══════════ Cliente HTTP ══════════ */

async function api(ruta, opciones = {}) {
  const cfg = { method: opciones.method || 'GET', headers: {} };
  if (token) cfg.headers['Authorization'] = 'Bearer ' + token;
  if (opciones.body !== undefined) {
    cfg.headers['Content-Type'] = 'application/json';
    cfg.body = JSON.stringify(opciones.body);
  }
  const res = await fetch(API + ruta, cfg);

  if (res.status === 401 && token) {
    cerrarSesion('Tu sesión expiró. Ingresa nuevamente.');
    throw new Error('sesión expirada');
  }
  if (res.status === 204) return null;

  const datos = await res.json().catch(() => null);
  if (!res.ok) {
    const mensaje = datos && datos.error ? datos.error : 'Error ' + res.status;
    throw new Error(mensaje);
  }
  return datos;
}

/* ══════════ Sesión y JWT ══════════ */

function decodificarJWT(t) {
  try {
    const b64 = t.split('.')[1].replace(/-/g, '+').replace(/_/g, '/');
    return JSON.parse(atob(b64));
  } catch {
    return null;
  }
}

// "recordar" decide dónde vive el token:
//  - sessionStorage (defecto): se borra al cerrar la pestaña — más seguro.
//  - localStorage (Recordarme): persiste entre aperturas de la app instalada,
//    con expiración validada en el cliente (exp) y en el servidor (401).
function guardarSesion(t, recordar) {
  const claims = decodificarJWT(t);
  if (!claims) throw new Error('token inválido');
  token = t;
  usuario = { email: claims.sub, rol: claims.rol, exp: claims.exp };
  (recordar ? localStorage : sessionStorage).setItem(CLAVE_TOKEN, t);
}

function restaurarSesion() {
  const t = sessionStorage.getItem(CLAVE_TOKEN) || localStorage.getItem(CLAVE_TOKEN);
  if (!t) return false;
  const claims = decodificarJWT(t);
  if (!claims || (claims.exp && claims.exp * 1000 < Date.now())) {
    sessionStorage.removeItem(CLAVE_TOKEN);
    localStorage.removeItem(CLAVE_TOKEN);
    return false;
  }
  token = t;
  usuario = { email: claims.sub, rol: claims.rol, exp: claims.exp };
  return true;
}

function cerrarSesion(mensaje) {
  token = null;
  usuario = null;
  sessionStorage.removeItem(CLAVE_TOKEN);
  localStorage.removeItem(CLAVE_TOKEN);
  mostrarLogin(mensaje);
}

const esAdmin = () => usuario && usuario.rol === 'admin';

/* ══════════ Vistas ══════════ */

function mostrarLogin(aviso) {
  $('#vista-app').classList.add('oculto');
  $('#vista-login').classList.remove('oculto');
  const nodoAviso = $('#login-aviso');
  if (aviso) {
    nodoAviso.textContent = aviso;
    nodoAviso.classList.remove('oculto');
  } else {
    nodoAviso.classList.add('oculto');
  }
}

function mostrarApp() {
  $('#vista-login').classList.add('oculto');
  $('#vista-app').classList.remove('oculto');
  $('#usuario-email').textContent = usuario.email;
  const rol = $('#usuario-rol');
  rol.textContent = usuario.rol;
  rol.className = 'badge ' + (esAdmin() ? 'badge-admin' : 'badge-veterinario');
  pintarPestanas();
  irAPestana('panel');
}

/* ══════════ Pestañas ══════════ */

const PESTANAS = [
  { id: 'panel',        icono: '📊', nombre: 'Panel' },
  { id: 'alertas',      icono: '🚨', nombre: 'Alertas' },
  { id: 'asignaciones', icono: '🩹', nombre: 'Triage' },
  { id: 'clinicas',     icono: '🏥', nombre: 'Clínicas' },
  { id: 'recursos',     icono: '⚙️', nombre: 'Recursos' },
  { id: 'mascotas',     icono: '🐶', nombre: 'Mascotas' },
  { id: 'historial',    icono: '📋', nombre: 'Historial' },
];

function pintarPestanas() {
  const nav = $('#pestanas');
  limpiar(nav);
  for (const p of PESTANAS) {
    const btn = el('button', 'pestana' + (p.id === pestanaActiva ? ' activa' : ''));
    btn.type = 'button';
    btn.appendChild(el('span', 'pestana-icono', p.icono));
    btn.appendChild(el('span', 'pestana-nombre', p.nombre));
    btn.addEventListener('click', () => irAPestana(p.id));
    nav.appendChild(btn);
  }
}

async function irAPestana(id) {
  pestanaActiva = id;
  pintarPestanas();
  const contenido = $('#contenido');
  limpiar(contenido);
  contenido.appendChild(el('p', 'descripcion', 'Cargando…'));
  try {
    await RENDER[id](contenido);
  } catch (err) {
    limpiar(contenido);
    const caja = el('div', 'tarjeta');
    caja.appendChild(el('h2', null, 'No se pudo cargar la sección'));
    caja.appendChild(el('p', 'descripcion', err.message));
    contenido.appendChild(caja);
  }
}

/* ══════════ Piezas reutilizables ══════════ */

function tarjeta(titulo, descripcion) {
  const caja = el('div', 'tarjeta');
  if (titulo) caja.appendChild(el('h2', null, titulo));
  if (descripcion) caja.appendChild(el('p', 'descripcion', descripcion));
  return caja;
}

function campo(labelTexto, input) {
  const caja = el('div', 'campo');
  const label = el('label', null, labelTexto);
  caja.appendChild(label);
  caja.appendChild(input);
  return caja;
}

function inputTexto(placeholder, requerido) {
  const i = el('input');
  i.type = 'text';
  i.placeholder = placeholder || '';
  if (requerido) i.required = true;
  return i;
}

function inputNumero(min, max, valor) {
  const i = el('input');
  i.type = 'number';
  if (min !== undefined) i.min = min;
  if (max !== undefined) i.max = max;
  if (valor !== undefined) i.value = valor;
  return i;
}

function selectDe(opciones, valorActual) {
  const s = el('select');
  for (const [valor, texto] of opciones) {
    const o = el('option', null, texto);
    o.value = valor;
    if (String(valor) === String(valorActual)) o.selected = true;
    s.appendChild(o);
  }
  return s;
}

function checkbox(labelTexto, marcado) {
  const caja = el('label', 'campo-check');
  const c = el('input');
  c.type = 'checkbox';
  c.checked = !!marcado;
  caja.appendChild(c);
  caja.appendChild(el('span', null, labelTexto));
  return { caja, input: c };
}

function botonEnviar(texto) {
  const b = el('button', 'btn btn-primario', texto);
  b.type = 'submit';
  return b;
}

// tabla(['ID', 'Nombre'], filas) — cada fila es un array de nodos o textos.
function tabla(cabeceras, filas, mensajeVacio) {
  const scroll = el('div', 'tabla-scroll');
  const t = el('table');
  const thead = el('thead');
  const trh = el('tr');
  for (const c of cabeceras) trh.appendChild(el('th', null, c));
  thead.appendChild(trh);
  t.appendChild(thead);

  const tbody = el('tbody');
  if (!filas.length) {
    const tr = el('tr');
    const td = el('td', 'vacio', mensajeVacio || 'Sin registros aún');
    td.colSpan = cabeceras.length;
    tr.appendChild(td);
    tbody.appendChild(tr);
  }
  for (const fila of filas) {
    const tr = el('tr');
    fila.forEach((celda, i) => {
      const td = el('td');
      td.dataset.label = cabeceras[i] || ''; // en móvil la tabla se ve como tarjetas
      if (celda instanceof Node) td.appendChild(celda);
      else td.textContent = celda === undefined || celda === null ? '—' : String(celda);
      tr.appendChild(td);
    });
    tbody.appendChild(tr);
  }
  t.appendChild(tbody);
  scroll.appendChild(t);
  return scroll;
}

function botonEliminar(ruta, recargar) {
  const b = el('button', 'btn btn-peligro btn-mini', 'Eliminar');
  b.type = 'button';
  b.addEventListener('click', async () => {
    if (!confirm('¿Eliminar este registro? Esta acción no se puede deshacer.')) return;
    try {
      await api(ruta, { method: 'DELETE' });
      toast('Registro eliminado', 'exito');
      recargar();
    } catch (err) {
      toast(err.message, 'error');
    }
  });
  return b;
}

function badge(texto, clase) {
  return el('span', 'badge ' + clase, texto);
}

function fechaLegible(iso) {
  if (!iso) return '—';
  const d = new Date(iso);
  return isNaN(d) ? iso : d.toLocaleString('es-EC', { dateStyle: 'short', timeStyle: 'short' });
}

const listar = async (ruta) => (await api(ruta)) || [];

/* ══════════ Pestaña: Panel ══════════ */

async function renderPanel(contenido) {
  const [alertas, mascotas, clinicas, recursos] = await Promise.all([
    listar('/alertas'), listar('/mascotas'), listar('/veterinarios'), listar('/recursos'),
  ]);
  limpiar(contenido);

  const activas = alertas.filter((a) => a.estado !== 'Atendido');
  const stats = el('div', 'estadisticas');
  const datosStats = [
    [activas.length, 'Alertas activas', 'rojo'],
    [mascotas.length, 'Mascotas registradas', 'azul'],
    [clinicas.filter((c) => c.activo).length, 'Clínicas activas', 'teal'],
    [recursos.filter((r) => r.esta_disponible).length, 'Recursos disponibles', 'verde'],
  ];
  for (const [valor, nombre, color] of datosStats) {
    const s = el('div', 'stat ' + color);
    s.appendChild(el('div', 'stat-valor', valor));
    s.appendChild(el('div', 'stat-nombre', nombre));
    stats.appendChild(s);
  }
  contenido.appendChild(stats);

  const criticas = alertas
    .filter((a) => a.gravedad >= 4 && a.estado !== 'Atendido')
    .sort((a, b) => b.gravedad - a.gravedad);

  const caja = tarjeta('🚨 Emergencias críticas', 'Alertas de gravedad 4-5 pendientes de atención.');
  const nombresMascota = new Map(mascotas.map((m) => [m.id, m.nombre]));
  caja.appendChild(tabla(
    ['Gravedad', 'Mascota', 'Requerimiento', 'Estado', 'Creada'],
    criticas.map((a) => [
      el('span', 'grav grav-' + a.gravedad, a.gravedad),
      nombresMascota.get(a.mascota_id) || 'Sin registrar',
      a.requerimiento,
      badge(a.estado, a.estado === 'Buscando' ? 'badge-amarillo' : 'badge-azul'),
      fechaLegible(a.creado_en),
    ]),
    'Sin emergencias críticas pendientes 🎉',
  ));
  contenido.appendChild(caja);
}

/* ══════════ Pestaña: Alertas ══════════ */

const ESTADOS_ALERTA = [['Buscando', 'Buscando'], ['Asignada', 'Asignada'], ['Atendido', 'Atendido']];

async function renderAlertas(contenido) {
  const [alertas, mascotas] = await Promise.all([listar('/alertas'), listar('/mascotas')]);
  limpiar(contenido);
  const recargar = () => irAPestana('alertas');
  const nombresMascota = new Map(mascotas.map((m) => [m.id, m.nombre]));

  // — Formulario de creación —
  const cajaForm = tarjeta('Nueva alerta de emergencia', 'La gravedad sigue la escala de triage: 1 (leve) a 5 (crítico).');
  const form = el('form', 'formulario');
  const selMascota = selectDe(
    [['0', '— Mascota sin registrar —']].concat(mascotas.map((m) => [m.id, `${m.nombre} (${m.especie})`])),
  );
  const selGravedad = selectDe([
    ['1', '1 — Leve'], ['2', '2 — Menor'], ['3', '3 — Moderada'], ['4', '4 — Grave'], ['5', '5 — Crítica'],
  ], '3');
  const inReq = inputTexto('Ej: Ventilador mecánico urgente', true);
  form.appendChild(campo('Mascota', selMascota));
  form.appendChild(campo('Gravedad', selGravedad));
  form.appendChild(campo('Requerimiento *', inReq));
  form.appendChild(botonEnviar('Crear alerta'));
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    try {
      await api('/alertas', { method: 'POST', body: {
        mascota_id: Number(selMascota.value),
        gravedad: Number(selGravedad.value),
        requerimiento: inReq.value.trim(),
        estado: 'Buscando',
      }});
      toast('Alerta creada', 'exito');
      recargar();
    } catch (err) { toast(err.message, 'error'); }
  });
  cajaForm.appendChild(form);
  contenido.appendChild(cajaForm);

  // — Listado —
  const cabeceras = ['ID', 'Gravedad', 'Mascota', 'Requerimiento', 'Estado', 'Creada'];
  if (esAdmin()) cabeceras.push('Acciones');

  const filas = alertas.map((a) => {
    const selEstado = selectDe(ESTADOS_ALERTA, a.estado);
    selEstado.addEventListener('change', async () => {
      try {
        await api('/alertas/' + a.id, { method: 'PUT', body: {
          mascota_id: a.mascota_id, gravedad: a.gravedad,
          requerimiento: a.requerimiento, estado: selEstado.value,
        }});
        toast('Estado actualizado', 'exito');
        recargar();
      } catch (err) { toast(err.message, 'error'); recargar(); }
    });
    const fila = [
      a.id,
      el('span', 'grav grav-' + a.gravedad, a.gravedad),
      nombresMascota.get(a.mascota_id) || 'Sin registrar',
      a.requerimiento,
      selEstado,
      fechaLegible(a.creado_en),
    ];
    if (esAdmin()) fila.push(botonEliminar('/alertas/' + a.id, recargar));
    return fila;
  });

  const cajaLista = tarjeta('Alertas registradas', 'Cambia el estado directamente desde la tabla.');
  cajaLista.appendChild(tabla(cabeceras, filas));
  contenido.appendChild(cajaLista);
}

/* ══════════ Pestaña: Asignaciones ══════════ */

const ESTADOS_ASIGNACION = [['Pendiente', 'Pendiente'], ['Confirmado', 'Confirmado'], ['Rechazado', 'Rechazado']];

async function renderAsignaciones(contenido) {
  const [asignaciones, alertas, recursos] = await Promise.all([
    listar('/asignaciones'), listar('/alertas'), listar('/recursos'),
  ]);
  limpiar(contenido);
  const recargar = () => irAPestana('asignaciones');

  const alertaTexto = new Map(alertas.map((a) => [a.id, `#${a.id} — ${a.requerimiento}`]));
  const recursoTexto = new Map(recursos.map((r) => [r.id, `#${r.id} — ${r.tipo_maquina}`]));

  const cajaForm = tarjeta('Asignar recurso a una alerta', 'Vincula un equipo clínico disponible con una emergencia activa.');
  const form = el('form', 'formulario');
  const selAlerta = selectDe(alertas.filter((a) => a.estado !== 'Atendido').map((a) => [a.id, alertaTexto.get(a.id)]));
  const selRecurso = selectDe(recursos.map((r) => [
    r.id, recursoTexto.get(r.id) + (r.esta_disponible ? '' : ' (no disponible)'),
  ]));
  form.appendChild(campo('Alerta *', selAlerta));
  form.appendChild(campo('Recurso *', selRecurso));
  form.appendChild(botonEnviar('Crear asignación'));
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    try {
      await api('/asignaciones', { method: 'POST', body: {
        alerta_id: Number(selAlerta.value),
        recurso_id: Number(selRecurso.value),
        estado_confirmacion: 'Pendiente',
      }});
      toast('Asignación creada', 'exito');
      recargar();
    } catch (err) { toast(err.message, 'error'); }
  });
  cajaForm.appendChild(form);
  contenido.appendChild(cajaForm);

  const cabeceras = ['ID', 'Alerta', 'Recurso', 'Confirmación'];
  if (esAdmin()) cabeceras.push('Acciones');

  const filas = asignaciones.map((asg) => {
    const selEstado = selectDe(ESTADOS_ASIGNACION, asg.estado_confirmacion);
    selEstado.addEventListener('change', async () => {
      try {
        await api('/asignaciones/' + asg.id, { method: 'PUT', body: {
          alerta_id: asg.alerta_id, recurso_id: asg.recurso_id,
          estado_confirmacion: selEstado.value,
        }});
        toast('Confirmación actualizada', 'exito');
        recargar();
      } catch (err) { toast(err.message, 'error'); recargar(); }
    });
    const fila = [
      asg.id,
      alertaTexto.get(asg.alerta_id) || '#' + asg.alerta_id,
      recursoTexto.get(asg.recurso_id) || '#' + asg.recurso_id,
      selEstado,
    ];
    if (esAdmin()) fila.push(botonEliminar('/asignaciones/' + asg.id, recargar));
    return fila;
  });

  const cajaLista = tarjeta('Asignaciones de triage');
  cajaLista.appendChild(tabla(cabeceras, filas));
  contenido.appendChild(cajaLista);
}

/* ══════════ Pestaña: Clínicas (perfiles veterinarios) ══════════ */

async function renderClinicas(contenido) {
  const clinicas = await listar('/veterinarios');
  limpiar(contenido);
  const recargar = () => irAPestana('clinicas');

  const cajaForm = tarjeta('Registrar clínica veterinaria', 'Nombre y teléfono son obligatorios.');
  const form = el('form', 'formulario');
  const inNombre = inputTexto('Ej: Clínica Los Ceibos', true);
  const inTelefono = inputTexto('Ej: +593 99 123 4567', true);
  inTelefono.type = 'tel';
  const inDireccion = inputTexto('Ej: Av. 4 de Noviembre, Manta');
  const chkActivo = checkbox('Activa', true);
  form.appendChild(campo('Nombre *', inNombre));
  form.appendChild(campo('Teléfono *', inTelefono));
  form.appendChild(campo('Dirección', inDireccion));
  form.appendChild(chkActivo.caja);
  form.appendChild(botonEnviar('Registrar clínica'));
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    try {
      await api('/veterinarios', { method: 'POST', body: {
        nombre: inNombre.value.trim(),
        telefono: inTelefono.value.trim(),
        direccion: inDireccion.value.trim(),
        activo: chkActivo.input.checked,
      }});
      toast('Clínica registrada', 'exito');
      recargar();
    } catch (err) { toast(err.message, 'error'); }
  });
  cajaForm.appendChild(form);
  contenido.appendChild(cajaForm);

  const cabeceras = ['ID', 'Nombre', 'Teléfono', 'Dirección', 'Estado', 'Equipos'];
  if (esAdmin()) cabeceras.push('Acciones');

  const filas = clinicas.map((c) => {
    const fila = [
      c.id, c.nombre, c.telefono, c.direccion || '—',
      badge(c.activo ? 'Activa' : 'Inactiva', c.activo ? 'badge-verde' : 'badge-gris'),
      (c.recursos || []).length,
    ];
    if (esAdmin()) fila.push(botonEliminar('/veterinarios/' + c.id, recargar));
    return fila;
  });

  const cajaLista = tarjeta('Clínicas registradas');
  cajaLista.appendChild(tabla(cabeceras, filas));
  contenido.appendChild(cajaLista);
}

/* ══════════ Pestaña: Recursos clínicos ══════════ */

async function renderRecursos(contenido) {
  const [recursos, clinicas] = await Promise.all([listar('/recursos'), listar('/veterinarios')]);
  limpiar(contenido);
  const recargar = () => irAPestana('recursos');
  const nombresClinica = new Map(clinicas.map((c) => [c.id, c.nombre]));

  const cajaForm = tarjeta('Registrar recurso clínico', 'Cada equipo pertenece a una clínica veterinaria.');
  const form = el('form', 'formulario');
  const selClinica = selectDe(clinicas.map((c) => [c.id, c.nombre]));
  const inTipo = inputTexto('Ej: Ecógrafo portátil', true);
  const chkDisponible = checkbox('Disponible', true);
  form.appendChild(campo('Clínica *', selClinica));
  form.appendChild(campo('Tipo de máquina *', inTipo));
  form.appendChild(chkDisponible.caja);
  form.appendChild(botonEnviar('Registrar recurso'));
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    if (!clinicas.length) { toast('Primero registra una clínica', 'error'); return; }
    try {
      await api('/recursos', { method: 'POST', body: {
        perfil_id: Number(selClinica.value),
        tipo_maquina: inTipo.value.trim(),
        esta_disponible: chkDisponible.input.checked,
      }});
      toast('Recurso registrado', 'exito');
      recargar();
    } catch (err) { toast(err.message, 'error'); }
  });
  cajaForm.appendChild(form);
  contenido.appendChild(cajaForm);

  const cabeceras = ['ID', 'Clínica', 'Tipo de máquina', 'Disponibilidad', 'Cambiar'];
  if (esAdmin()) cabeceras.push('Acciones');

  const filas = recursos.map((r) => {
    const btnToggle = el('button', 'btn btn-suave btn-mini', r.esta_disponible ? 'Marcar ocupado' : 'Marcar libre');
    btnToggle.type = 'button';
    btnToggle.addEventListener('click', async () => {
      try {
        await api('/recursos/' + r.id, { method: 'PUT', body: {
          perfil_id: r.perfil_id, tipo_maquina: r.tipo_maquina,
          esta_disponible: !r.esta_disponible,
        }});
        toast('Disponibilidad actualizada', 'exito');
        recargar();
      } catch (err) { toast(err.message, 'error'); }
    });
    const fila = [
      r.id,
      nombresClinica.get(r.perfil_id) || '#' + r.perfil_id,
      r.tipo_maquina,
      badge(r.esta_disponible ? 'Disponible' : 'Ocupado', r.esta_disponible ? 'badge-verde' : 'badge-rojo'),
      btnToggle,
    ];
    if (esAdmin()) fila.push(botonEliminar('/recursos/' + r.id, recargar));
    return fila;
  });

  const cajaLista = tarjeta('Inventario de equipos');
  cajaLista.appendChild(tabla(cabeceras, filas));
  contenido.appendChild(cajaLista);
}

/* ══════════ Pestaña: Mascotas ══════════ */

async function renderMascotas(contenido) {
  const mascotas = await listar('/mascotas');
  limpiar(contenido);
  const recargar = () => irAPestana('mascotas');

  const cajaForm = tarjeta('Registrar mascota', 'El nombre es obligatorio.');
  const form = el('form', 'formulario');
  const inNombre = inputTexto('Ej: Rex', true);
  const inEspecie = inputTexto('Ej: Perro', true);
  const inEdad = inputNumero(0, 60, 1);
  const inDueno = inputTexto('Ej: Ana Pérez', true);
  form.appendChild(campo('Nombre *', inNombre));
  form.appendChild(campo('Especie *', inEspecie));
  form.appendChild(campo('Edad (años)', inEdad));
  form.appendChild(campo('Dueño *', inDueno));
  form.appendChild(botonEnviar('Registrar mascota'));
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    try {
      await api('/mascotas', { method: 'POST', body: {
        nombre: inNombre.value.trim(),
        especie: inEspecie.value.trim(),
        edad: Number(inEdad.value) || 0,
        dueno: inDueno.value.trim(),
      }});
      toast('Mascota registrada', 'exito');
      recargar();
    } catch (err) { toast(err.message, 'error'); }
  });
  cajaForm.appendChild(form);
  contenido.appendChild(cajaForm);

  const cabeceras = ['ID', 'Nombre', 'Especie', 'Edad', 'Dueño', 'Entradas de historial'];
  if (esAdmin()) cabeceras.push('Acciones');

  const filas = mascotas.map((m) => {
    const fila = [m.id, m.nombre, m.especie, m.edad + ' años', m.dueno, (m.historial || []).length];
    if (esAdmin()) fila.push(botonEliminar('/mascotas/' + m.id, recargar));
    return fila;
  });

  const cajaLista = tarjeta('Pacientes registrados');
  cajaLista.appendChild(tabla(cabeceras, filas));
  contenido.appendChild(cajaLista);
}

/* ══════════ Pestaña: Historial médico ══════════ */

async function renderHistorial(contenido) {
  const [historial, mascotas] = await Promise.all([listar('/historial'), listar('/mascotas')]);
  limpiar(contenido);
  const recargar = () => irAPestana('historial');
  const nombresMascota = new Map(mascotas.map((m) => [m.id, `${m.nombre} (${m.especie})`]));

  const cajaForm = tarjeta('Nueva entrada clínica', 'Cada entrada pertenece a una mascota registrada.');
  const form = el('form', 'formulario');
  const selMascota = selectDe(mascotas.map((m) => [m.id, nombresMascota.get(m.id)]));
  const inDiagnostico = inputTexto('Ej: Gastritis leve', true);
  const inTratamiento = inputTexto('Ej: Dieta blanda 5 días');
  const inFecha = el('input');
  inFecha.type = 'date';
  inFecha.required = true;
  inFecha.value = new Date().toISOString().slice(0, 10);
  const inVeterinario = inputTexto('Ej: Dra. Vélez');
  form.appendChild(campo('Mascota *', selMascota));
  form.appendChild(campo('Diagnóstico *', inDiagnostico));
  form.appendChild(campo('Tratamiento', inTratamiento));
  form.appendChild(campo('Fecha *', inFecha));
  form.appendChild(campo('Veterinario', inVeterinario));
  form.appendChild(botonEnviar('Guardar entrada'));
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    if (!mascotas.length) { toast('Primero registra una mascota', 'error'); return; }
    try {
      await api('/historial', { method: 'POST', body: {
        mascota_id: Number(selMascota.value),
        diagnostico: inDiagnostico.value.trim(),
        tratamiento: inTratamiento.value.trim(),
        fecha: inFecha.value,
        veterinario: inVeterinario.value.trim(),
      }});
      toast('Entrada guardada', 'exito');
      recargar();
    } catch (err) { toast(err.message, 'error'); }
  });
  cajaForm.appendChild(form);
  contenido.appendChild(cajaForm);

  const filas = historial.map((h) => [
    h.id,
    nombresMascota.get(h.mascota_id) || '#' + h.mascota_id,
    h.fecha,
    h.diagnostico,
    h.tratamiento || '—',
    h.veterinario || '—',
  ]);

  const cajaLista = tarjeta('Historial médico completo');
  cajaLista.appendChild(tabla(['ID', 'Mascota', 'Fecha', 'Diagnóstico', 'Tratamiento', 'Veterinario'], filas));
  contenido.appendChild(cajaLista);
}

const RENDER = {
  panel: renderPanel,
  alertas: renderAlertas,
  asignaciones: renderAsignaciones,
  clinicas: renderClinicas,
  recursos: renderRecursos,
  mascotas: renderMascotas,
  historial: renderHistorial,
};

/* ══════════ Login / Registro ══════════ */

let modoRegistro = false;

function alternarModo() {
  modoRegistro = !modoRegistro;
  $('#btn-login').textContent = modoRegistro ? 'Crear cuenta' : 'Ingresar';
  $('#texto-alternar').textContent = modoRegistro ? '¿Ya tienes cuenta?' : '¿No tienes cuenta?';
  $('#btn-alternar').textContent = modoRegistro ? 'Inicia sesión' : 'Regístrate';
  $('#login-password').autocomplete = modoRegistro ? 'new-password' : 'current-password';
  $('#login-password').minLength = modoRegistro ? 8 : 6;
}

async function enviarLogin(e) {
  e.preventDefault();
  const btn = $('#btn-login');
  const email = $('#login-email').value.trim();
  const password = $('#login-password').value;
  btn.disabled = true;
  try {
    if (modoRegistro) {
      await api('/auth/register', { method: 'POST', body: { email, password } });
      toast('Cuenta creada, iniciando sesión…', 'exito');
    }
    const resp = await api('/auth/login', { method: 'POST', body: { email, password } });
    guardarSesion(resp.token, $('#login-recordar').checked);
    $('#login-password').value = '';
    mostrarApp();
  } catch (err) {
    const aviso = $('#login-aviso');
    aviso.textContent = err.message;
    aviso.className = 'aviso error';
  } finally {
    btn.disabled = false;
  }
}

function rellenarDemo(email, password) {
  $('#login-email').value = email;
  $('#login-password').value = password;
  if (modoRegistro) alternarModo();
}

/* ══════════ PWA: service worker e instalación ══════════ */

function registrarPWA() {
  if ('serviceWorker' in navigator) {
    navigator.serviceWorker.register('sw.js').catch(() => { /* opcional: la app funciona igual sin SW */ });
  }

  let eventoInstalar = null;
  const botones = document.querySelectorAll('.btn-instalar');

  window.addEventListener('beforeinstallprompt', (e) => {
    e.preventDefault();
    eventoInstalar = e;
    botones.forEach((b) => b.classList.remove('oculto'));
  });

  botones.forEach((b) => b.addEventListener('click', async () => {
    if (!eventoInstalar) return;
    eventoInstalar.prompt();
    await eventoInstalar.userChoice;
    eventoInstalar = null;
    botones.forEach((x) => x.classList.add('oculto'));
  }));

  window.addEventListener('appinstalled', () => {
    botones.forEach((b) => b.classList.add('oculto'));
    toast('Aplicación instalada 🎉', 'exito');
  });
}

/* ══════════ Arranque ══════════ */

document.addEventListener('DOMContentLoaded', () => {
  registrarPWA();
  $('#form-login').addEventListener('submit', enviarLogin);
  $('#btn-alternar').addEventListener('click', alternarModo);
  $('#btn-salir').addEventListener('click', () => cerrarSesion());
  $('#btn-demo-admin').addEventListener('click', () => rellenarDemo('admin@codigogarra.vet', 'Admin123!'));
  $('#btn-demo-vet').addEventListener('click', () => rellenarDemo('vet@codigogarra.vet', 'Vet123!'));

  if (restaurarSesion()) mostrarApp();
  else mostrarLogin();
});
