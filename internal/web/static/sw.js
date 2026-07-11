'use strict';

/* Service worker de Código Garra.
   Estrategia: red primero con respaldo en caché (app shell) para que la
   aplicación instalada abra incluso sin conexión.
   Seguridad: las peticiones a /api/ NUNCA se interceptan ni se cachean —
   los datos autenticados jamás quedan almacenados en el dispositivo. */

const CACHE = 'garra-shell-v1';
const SHELL = [
  '/',
  '/index.html',
  '/styles.css',
  '/app.js',
  '/manifest.webmanifest',
  '/icon-192.png',
  '/icon-512.png',
  '/icon-maskable-512.png',
];

self.addEventListener('install', (evento) => {
  evento.waitUntil(
    caches.open(CACHE)
      .then((cache) => cache.addAll(SHELL))
      .then(() => self.skipWaiting()),
  );
});

self.addEventListener('activate', (evento) => {
  evento.waitUntil(
    caches.keys()
      .then((claves) => Promise.all(claves.filter((c) => c !== CACHE).map((c) => caches.delete(c))))
      .then(() => self.clients.claim()),
  );
});

self.addEventListener('fetch', (evento) => {
  const url = new URL(evento.request.url);

  // La API se va siempre a la red: sin caché de datos sensibles.
  if (url.pathname.startsWith('/api/')) return;
  if (evento.request.method !== 'GET') return;

  evento.respondWith(
    fetch(evento.request)
      .then((respuesta) => {
        const copia = respuesta.clone();
        caches.open(CACHE).then((cache) => cache.put(evento.request, copia));
        return respuesta;
      })
      .catch(() =>
        caches.match(evento.request, { ignoreSearch: true })
          .then((guardada) => guardada || caches.match('/')),
      ),
  );
});
