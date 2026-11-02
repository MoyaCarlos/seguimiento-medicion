# Quickstart: Markdown seguro en el Product Backlog

## 1. Instalar dependencias y levantar el frontend

```sh
cd web
npm install
npm run dev
```

## 2. Verificar el renderizado

1. Crear o editar una historia con una descripción como:

   ```markdown
   **Importante:** validar el alta.

   - primer paso
   - segundo paso

   ```go
   fmt.Println("ok")
   ```
   ```

2. Guardar y abrir `/projects/<id>/backlog`.
3. Comprobar que la descripción muestra negrita, lista y bloque de código, y que los
   criterios de aceptación también se interpretan.
4. Abrir la edición: los campos deben conservar el Markdown original.
5. Guardar sin cambios: el cuerpo de `PUT /backlog/{id}` debe llevar el mismo texto.

## 3. Verificar que no hay XSS

1. Editar la descripción y pegar:

   ```html
   <script>alert('xss')</script>
   <img src=x onerror="alert('xss')">
   ```

2. Guardar y volver al listado.
3. Abrir las herramientas de desarrollo: no debe existir un `<script>` inyectado ni un
   `<img>` con `onerror`; no debe aparecer ningún diálogo.

## 4. Automatizado

```sh
cd web
npm test
npm run lint
npm run build
```

`Markdown.test.tsx` cubre CommonMark y los intentos de XSS; `Listado.test.tsx` cubre la
integración en el listado.
