/// <reference types="vite/client" />
// This file enables typescript types for import.meta.env

interface ImportMetaEnv {
  readonly VITE_P2P_API_ROOT: string
  readonly VITE_P2P_TOKEN: string
}
