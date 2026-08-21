import { afterEach } from 'vitest'
import { cleanup } from '@testing-library/react'
import '@testing-library/jest-dom/vitest'

class MemoryStorage {
  private map = new Map<string, string>()
  getItem(k: string) {
    return this.map.has(k) ? this.map.get(k)! : null
  }
  setItem(k: string, v: string) {
    this.map.set(k, String(v))
  }
  removeItem(k: string) {
    this.map.delete(k)
  }
  clear() {
    this.map.clear()
  }
  key(i: number) {
    return Array.from(this.map.keys())[i] ?? null
  }
  get length() {
    return this.map.size
  }
}

if (typeof globalThis.localStorage === 'undefined') {
  const store = new MemoryStorage()
  Object.defineProperty(globalThis, 'localStorage', { value: store, configurable: true })
  Object.defineProperty(globalThis, 'sessionStorage', { value: new MemoryStorage(), configurable: true })
}

afterEach(() => {
  cleanup()
})
