const TRAZOS = {
  marca: 'M6 19v-6M12 19V5M18 19v-9',
  carpeta: 'M3 6h6l2 2h10v11H3z',
  panel: 'M4 4h7v9H4zM13 4h7v5h-7zM13 11h7v9h-7zM4 15h7v5H4z',
  lista: 'M9 6h11M9 12h11M9 18h11M4 6h1M4 12h1M4 18h1',
  bandera: 'M5 21V4h13l-3 4 3 4H5',
  defecto: 'M12 8a4 4 0 0 1 4 4v3a4 4 0 0 1-8 0v-3a4 4 0 0 1 4-4zM9.5 5l1 3M14.5 5l-1 3M4 13h4M16 13h4M5 19l3-2M19 19l-3-2',
  estado: 'M4 4v16h16M8 16v-4M12 16V8M16 16v-6',
  barras: 'M3 20h18M6 20v-8M12 20V5M18 20v-11',
  historial: 'M12 7v5l3 2M3 12a9 9 0 1 0 9-9',
  menu: 'M4 6h16M4 12h16M4 18h16',
  cerrar: 'M6 6l12 12M18 6L6 18',
  salir: 'M10 4H5v16h5M15 8l4 4-4 4M19 12H9',
  mas: 'M12 5v14M5 12h14',
  equipo: 'M16 20v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2M9 10a4 4 0 1 0 0-8 4 4 0 0 0 0 8zM22 20v-2a4 4 0 0 0-3-3.9M16 2.1a4 4 0 0 1 0 7.8',
  editar: 'M12 20h9M16.5 3.5a2.1 2.1 0 0 1 3 3L7 19l-4 1 1-4z',
  chevronAbajo: 'm6 9 6 6 6-6',
  calendario: 'M5 4h14a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2zM16 2v4M8 2v4M3 10h18',
  ojo: 'M2 12s3.6-7 10-7 10 7 10 7-3.6 7-10 7S2 12 2 12zM12 15a3 3 0 1 0 0-6 3 3 0 0 0 0 6z',
  ojoTachado: 'M17.9 17.9A10.9 10.9 0 0 1 12 19c-6.4 0-10-7-10-7a18.5 18.5 0 0 1 5.1-5.9M9.9 5.2A9.9 9.9 0 0 1 12 5c6.4 0 10 7 10 7a18.6 18.6 0 0 1-2.2 3.2M14.1 14.1a3 3 0 1 1-4.2-4.2M1 1l22 22',
}

export type NombreIcono = keyof typeof TRAZOS

export function Icono({ nombre }: { nombre: NombreIcono }) {
  return (
    <svg
      className="icono"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
      focusable="false"
    >
      <path d={TRAZOS[nombre]} />
    </svg>
  )
}
