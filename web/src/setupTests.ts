import { cleanup, configure } from '@testing-library/react'
import { afterEach } from 'vitest'

// Con la máquina cargada, una pantalla diferida puede tardar varios segundos en importarse.
configure({ asyncUtilTimeout: 10_000 })

afterEach(cleanup)
