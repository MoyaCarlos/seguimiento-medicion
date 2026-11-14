import { afterEach, beforeEach, describe, expect, it } from 'vitest'
import { porcentajeTranscurrido } from './avance'

const local = (anio: number, mes: number, dia: number, hora = 0) => new Date(anio, mes - 1, dia, hora)

describe('porcentajeTranscurrido', () => {
  it('Given un día anterior al inicio, When se calcula, Then es 0', () => {
    expect(porcentajeTranscurrido('2026-10-10', '2026-10-30', local(2026, 10, 1))).toBe(0)
  })

  it('Given el mismo día del inicio, When se calcula, Then es 0', () => {
    expect(porcentajeTranscurrido('2026-10-10', '2026-10-30', local(2026, 10, 10))).toBe(0)
  })

  it('Given la mitad del período, When se calcula, Then es 50', () => {
    expect(porcentajeTranscurrido('2026-10-10', '2026-10-30', local(2026, 10, 20))).toBe(50)
  })

  it('Given un día intermedio, When se calcula, Then se redondea al entero más cercano', () => {
    expect(porcentajeTranscurrido('2026-10-01', '2026-10-31', local(2026, 10, 11))).toBe(33)
  })

  it('Given dos tercios del período, When se calcula, Then redondea hacia arriba a 67 (no trunca)', () => {
    expect(porcentajeTranscurrido('2026-10-01', '2026-10-31', local(2026, 10, 21))).toBe(67)
  })

  it('Given el día de fin o uno posterior, When se calcula, Then es 100', () => {
    expect(porcentajeTranscurrido('2026-10-10', '2026-10-30', local(2026, 10, 30))).toBe(100)
    expect(porcentajeTranscurrido('2026-10-10', '2026-10-30', local(2026, 12, 1))).toBe(100)
  })

  it('Given un proyecto sin fecha de inicio o sin fecha de fin, When se calcula, Then no hay porcentaje', () => {
    expect(porcentajeTranscurrido(null, '2026-10-30', local(2026, 10, 20))).toBeNull()
    expect(porcentajeTranscurrido('2026-10-10', null, local(2026, 10, 20))).toBeNull()
  })

  describe('con una zona horaria al oeste de UTC', () => {
    const original = process.env.TZ
    beforeEach(() => {
      process.env.TZ = 'America/Argentina/Buenos_Aires'
    })
    afterEach(() => {
      if (original === undefined) delete process.env.TZ
      else process.env.TZ = original
    })

    it('Given la medianoche local del día de inicio, When se calcula, Then las fechas se leen como locales y da 0 (no como UTC)', () => {
      expect(porcentajeTranscurrido('2026-10-10', '2026-10-12', local(2026, 10, 10))).toBe(0)
    })

    it('Given la medianoche local del día de fin, When se calcula, Then da 100', () => {
      expect(porcentajeTranscurrido('2026-10-10', '2026-10-12', local(2026, 10, 12))).toBe(100)
    })

    it('Given una hora tarde del día anterior al fin, When se calcula, Then todavía no llegó a 100', () => {
      expect(porcentajeTranscurrido('2026-10-10', '2026-10-12', local(2026, 10, 11, 23))).toBeLessThan(100)
    })
  })
})
