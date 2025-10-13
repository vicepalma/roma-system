export const Bus = {
  emit<T = any>(name: string, detail?: T) {
    window.dispatchEvent(new CustomEvent(name, { detail }))
  },
  on<T = any>(name: string, handler: (detail: T) => void) {
    const listener = (e: Event) => handler((e as CustomEvent).detail as T)
    window.addEventListener(name, listener)
    return () => window.removeEventListener(name, listener)
  }
}
