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
