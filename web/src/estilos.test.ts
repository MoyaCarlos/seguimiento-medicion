/// <reference types="node" />
import { readFileSync } from 'node:fs'
import { expect, it } from 'vitest'

const estilos = readFileSync('src/App.css', 'utf8')

it('Given un enlace con clase button, When se estiliza, Then no se subraya igual que un botón', () => {
  const regla = estilos.match(/^\.button\s*\{[^}]*\}/m)?.[0] ?? ''

  expect(regla).toMatch(/text-decoration:\s*none/)
})

it('Given una pantalla de 320 px con barra de scroll que ocupa lugar, When la página tiene scroll vertical, Then body no fija un ancho mínimo que genere scroll horizontal', () => {
  const regla = estilos.match(/^body\s*\{[^}]*\}/m)?.[0] ?? ''

  expect(regla).not.toBe('')
  expect(regla).not.toMatch(/min-width/)
})

it('Given el mensaje de datos de prueba cargados, When se estiliza, Then es una caja verde suave con texto que cumple 4,5:1', () => {
  const css = readFileSync('src/features/autenticacion/estilos.css', 'utf8')
  const regla = css.match(/^\.demo-exito\s*\{([^}]*)\}/m)?.[1] ?? ''
  const fondo = regla.match(/background:\s*(#[0-9a-f]{6})/)?.[1] ?? ''
  const texto = regla.match(/\bcolor:\s*(#[0-9a-f]{6})/)?.[1] ?? ''
  const luz = (h: string) =>
    [1, 3, 5].map((i) => parseInt(h.slice(i, i + 2), 16) / 255).map((c) => (c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4))
      .reduce((a, c, i) => a + c * [0.2126, 0.7152, 0.0722][i], 0)

  expect(fondo).toBe('#dcf5e4')
  const [a, b] = [luz(fondo), luz(texto)].sort((x, y) => y - x)
  expect((a + 0.05) / (b + 0.05)).toBeGreaterThanOrEqual(4.5)
})

const bloque = (selector: string) =>
  estilos.match(new RegExp(`^${selector.replace(/[.:[\]]/g, '\\$&')}\\s*\\{([^}]*)\\}`, 'm'))?.[1] ?? ''

it('Given el encabezado de una pantalla, When se estiliza globalmente, Then h1 no lleva margen, el pretítulo queda a 6px y la línea del proyecto a 6px del título', () => {
  expect(bloque('h1')).toMatch(/margin:\s*0\s*;/)
  expect(bloque('.eyebrow')).toMatch(/margin:\s*0 0 6px/)
  expect(bloque('.project-context')).toMatch(/margin:\s*6px 0 0/)
  expect(bloque('.project-context')).toMatch(/font-size:\s*14\.5px/)
  expect(estilos).not.toMatch(/\.page-header \.eyebrow\s*\{/)
})

it('Given Estado y Dashboard, When el encabezado ya tiene su separación global, Then no compensan el margen del pretítulo con márgenes negativos', () => {
  const estado = readFileSync('src/features/estado/estado.css', 'utf8')

  expect(estado).not.toMatch(/margin-bottom:\s*-/)
})

it('Given el ítem activo de la barra lateral, When se estiliza, Then la barra de acento va pegada al borde de la barra lateral y no dentro del ítem', () => {
  const activo = bloque('.barra-lateral__enlace--activo')
  const barra = bloque('.barra-lateral__enlace--activo::before')

  expect(activo).not.toMatch(/box-shadow:\s*inset/)
  expect(barra).toMatch(/left:\s*-14px/)
  expect(barra).toMatch(/width:\s*3px/)
  expect(barra).toMatch(/background:\s*var\(--accent\)/)
  expect(bloque('.barra-lateral__enlace')).toMatch(/position:\s*relative/)
})

it('Given las pestañas inferiores, When se estilizan, Then van fijas abajo, no tapan el final del contenido y su texto cumple 4,5:1', () => {
  const angosto = estilos.slice(estilos.indexOf('@media (max-width: 900px)'))
  const pestanas = angosto.match(/\.pestanas\s*\{([^}]*)\}/)?.[1] ?? ''
  const activa = angosto.match(/\.pestanas__enlace--activa\s*\{([^}]*)\}/)?.[1] ?? ''
  const contenido = angosto.match(/\.armazon--con-pestanas \.contenido\s*\{([^}]*)\}/)?.[1] ?? ''
  const luz = (h: string) =>
    [1, 3, 5].map((i) => parseInt(h.slice(i, i + 2), 16) / 255).map((c) => (c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4))
      .reduce((a, c, i) => a + c * [0.2126, 0.7152, 0.0722][i], 0)
  const contraste = (x: string, y: string) => {
    const [a, b] = [luz(x), luz(y)].sort((p, q) => q - p)
    return (a + 0.05) / (b + 0.05)
  }
  const fondo = activa.match(/background:\s*(#[0-9a-f]{6})/)?.[1] ?? ''
  const texto = activa.match(/\bcolor:\s*(#[0-9a-f]{6})/)?.[1] ?? ''

  expect(pestanas).toMatch(/position:\s*fixed/)
  expect(pestanas).toMatch(/bottom:\s*0/)
  expect(Number(contenido.match(/padding-bottom:\s*(\d+)px/)?.[1])).toBeGreaterThanOrEqual(96)
  expect(contraste(texto, fondo)).toBeGreaterThanOrEqual(4.5)
  expect(contraste('#606f86', '#ffffff')).toBeGreaterThanOrEqual(4.5)
})

it('Given el selector de proyecto, When se estiliza, Then el resaltado del teclado es un contorno de 2px que se distingue 3:1 del fondo y las iniciales cumplen 4,5:1', () => {
  const resaltada = bloque('.selector-proyecto__lista:focus-visible .selector-proyecto__opcion--resaltada')
  const luz = (h: string) =>
    [1, 3, 5].map((i) => parseInt(h.slice(i, i + 2), 16) / 255).map((c) => (c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4))
      .reduce((a, c, i) => a + c * [0.2126, 0.7152, 0.0722][i], 0)
  const contraste = (x: string, y: string) => {
    const [a, b] = [luz(x), luz(y)].sort((p, q) => q - p)
    return (a + 0.05) / (b + 0.05)
  }

  expect(resaltada).toMatch(/box-shadow:\s*inset 0 0 0 2px var\(--brand\)/)
  expect(contraste('#0d9488', '#ffffff')).toBeGreaterThanOrEqual(3)
  for (const n of [0, 1, 2, 3, 4]) {
    const avatar = bloque(`.avatar-proyecto--${n}`)
    const fondo = avatar.match(/background:\s*(#[0-9a-f]{6})/)?.[1] ?? ''
    const texto = avatar.match(/\bcolor:\s*(#[0-9a-f]{6})/)?.[1] ?? ''
    expect(contraste(texto, fondo)).toBeGreaterThanOrEqual(4.5)
  }
  expect(bloque('.selector-proyecto__popup[hidden]')).toMatch(/display:\s*none/)
})

it('Given el texto en cursiva del prototipo, When se estiliza, Then el navegador puede dibujarla (Inter no trae cursiva propia)', () => {
  expect(estilos).not.toMatch(/font-synthesis:\s*none/)
  expect(estilos).toMatch(/font-synthesis:\s*style/)
})

it('Given la ruta de navegación, When se estiliza, Then las migas truncan a 24ch con puntos suspensivos y se ven también en pantallas angostas', () => {
  const miga = estilos.match(/^\.ruta a,\s*\.ruta \[aria-current\]\s*\{([^}]*)\}/m)?.[1] ?? ''

  expect(miga).toMatch(/max-width:\s*24ch/)
  expect(miga).toMatch(/text-overflow:\s*ellipsis/)
  expect(estilos).not.toMatch(/\.ruta > :not\(\[aria-current\]\)/)
})

it('Given los separadores de la ruta, When se estilizan, Then son iconos de 14px algo transparentes', () => {
  const separador = estilos.match(/^\.ruta svg\s*\{([^}]*)\}/m)?.[1] ?? ''

  expect(separador).toMatch(/height:\s*14px/)
  expect(separador).toMatch(/width:\s*14px/)
  expect(separador).toMatch(/opacity:\s*0?\.6/)
})

it('Given un botón discreto gris (ghost), When se estiliza, Then es transparente, gris y cumple 4,5:1 contra blanco', () => {
  const cuerpo = bloque('.button--ghost')
  const texto = cuerpo.match(/color:\s*(#[0-9a-f]{6})/)?.[1] ?? ''
  const luz = (h: string) =>
    [1, 3, 5].map((i) => parseInt(h.slice(i, i + 2), 16) / 255).map((c) => (c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4))
      .reduce((a, c, i) => a + c * [0.2126, 0.7152, 0.0722][i], 0)

  expect(cuerpo).toMatch(/background:\s*transparent/)
  expect(texto).not.toBe('')
  expect((1.05) / (luz(texto) + 0.05)).toBeGreaterThanOrEqual(4.5)
  expect(bloque('.button--ghost:hover')).toMatch(/background:\s*#eef2f6/)
})

it('Given la sección «Proyectos cerrados», When se estiliza, Then su título es una etiqueta de 12px en mayúsculas separada de los abiertos', () => {
  expect(bloque('.lista-proyectos__cerrados h2')).toMatch(/font-size:\s*12px/)
  expect(bloque('.lista-proyectos__cerrados h2')).toMatch(/text-transform:\s*uppercase/)
  expect(bloque('.lista-proyectos__cerrados')).toMatch(/gap:\s*12px/)
})

it('Given la insignia «Cerrado el …», When se estiliza, Then es una píldora gris de 12px y cumple 4,5:1', () => {
  const cuerpo = bloque('.tarjeta-proyecto__cierre')
  const texto = cuerpo.match(/color:\s*(#[0-9a-f]{6})/)?.[1] ?? ''
  const fondo = cuerpo.match(/background:\s*(#[0-9a-f]{6})/)?.[1] ?? ''
  const luz = (h: string) =>
    [1, 3, 5].map((i) => parseInt(h.slice(i, i + 2), 16) / 255).map((c) => (c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4))
      .reduce((a, c, i) => a + c * [0.2126, 0.7152, 0.0722][i], 0)

  expect(cuerpo).toMatch(/border-radius:\s*999px/)
  expect(cuerpo).toMatch(/font-size:\s*12px/)
  expect(texto).not.toBe('')
  const [a, b] = [luz(fondo), luz(texto)].sort((x, y) => y - x)
  expect((a + 0.05) / (b + 0.05)).toBeGreaterThanOrEqual(4.5)
})

it('Given las etiquetas de la tarjeta de proyecto, When se estilizan, Then van a la derecha con 6px entre sí y el rol es una etiqueta gris de 12px que cumple 4,5:1', () => {
  const etiquetas = bloque('.tarjeta-proyecto__etiquetas')
  const rol = bloque('.tarjeta-proyecto__rol')
  const texto = rol.match(/color:\s*(#[0-9a-f]{6})/)?.[1] ?? ''
  const fondo = rol.match(/background:\s*(#[0-9a-f]{6})/)?.[1] ?? ''
  const luz = (h: string) =>
    [1, 3, 5].map((i) => parseInt(h.slice(i, i + 2), 16) / 255).map((c) => (c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4))
      .reduce((a, c, i) => a + c * [0.2126, 0.7152, 0.0722][i], 0)

  expect(etiquetas).toMatch(/margin-left:\s*auto/)
  expect(etiquetas).toMatch(/gap:\s*6px/)
  expect(rol).toMatch(/border-radius:\s*7px/)
  expect(rol).toMatch(/font-size:\s*12px/)
  const [a, b] = [luz(fondo), luz(texto)].sort((x, y) => y - x)
  expect((a + 0.05) / (b + 0.05)).toBeGreaterThanOrEqual(4.5)
})

it('Given la sección de cerrados, When se estiliza, Then la separa una línea divisoria de los abiertos y el título lleva su contador', () => {
  const seccion = bloque('.lista-proyectos__cerrados')

  expect(seccion).toMatch(/border-top:\s*1px solid var\(--line-strong\)/)
  expect(seccion).toMatch(/padding-top:\s*20px/)
  expect(seccion).toMatch(/margin-top:\s*16px/)
  expect(bloque('.lista-proyectos__cerrados-cabecera')).toMatch(/display:\s*flex/)
  expect(bloque('.lista-proyectos__contador')).toMatch(/border-radius:\s*999px/)
})

it('Given las tarjetas de los cerrados, When se estilizan, Then quedan más tenues que las de los abiertos (fondo gris suave y sin sombra) sin perder contraste de texto', () => {
  const tarjeta = bloque('.lista-proyectos__cerrados .tarjeta-proyecto')
  const fondo = tarjeta.match(/background:\s*(#[0-9a-f]{6})/)?.[1] ?? ''
  const luz = (h: string) =>
    [1, 3, 5].map((i) => parseInt(h.slice(i, i + 2), 16) / 255).map((c) => (c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4))
      .reduce((a, c, i) => a + c * [0.2126, 0.7152, 0.0722][i], 0)

  expect(fondo).toBe('#f8fafc')
  expect(tarjeta).toMatch(/box-shadow:\s*none/)
  const [claro, oscuro] = [luz(fondo), luz('#606f86')].sort((x, y) => y - x)
  expect((claro + 0.05) / (oscuro + 0.05)).toBeGreaterThanOrEqual(4.5)
})
