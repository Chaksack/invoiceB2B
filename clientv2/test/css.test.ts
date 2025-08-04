import { describe, it, expect } from 'vitest'
import '~/assets/css/tailwind.css'

describe('CSS Import Test', () => {
  it('should import tailwind.css without errors', () => {
    // If the import above doesn't throw an error, the test passes
    expect(true).toBe(true)
  })
})