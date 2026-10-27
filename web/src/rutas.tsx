import { Navigate, type RouteObject } from 'react-router'
import App from './App'
import { LayoutProyecto } from './compartido/router/LayoutProyecto'
import { NoEncontrada } from './compartido/router/NoEncontrada'

export const rutas: RouteObject[] = [
  {
    element: <App />,
    children: [
      { index: true, lazy: () => import('./features/proyecto/crear') },
      {
        path: 'projects/:id',
        element: <LayoutProyecto />,
        children: [
          { index: true, element: <Navigate to="proyecto" replace /> },
          { path: 'proyecto/*', lazy: () => import('./features/proyecto') },
          { path: 'backlog/*', lazy: () => import('./features/backlog') },
          { path: 'sprints/*', lazy: () => import('./features/sprints') },
          { path: 'status/*', lazy: () => import('./features/estado') },
        ],
      },
      { path: '*', element: <NoEncontrada /> },
    ],
  },
]
