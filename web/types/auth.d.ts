import type { PbUser } from "~/src/client"

declare module '#auth-utils' {
  interface User extends PbUser { }

  interface UserSession {
    id: string
  }

  interface SecureSessionData {
    access_token: string,
    refresh_token: string
  }


  interface UserSessionComposable {
    /**
     * Computed indicating if the auth session is ready
     */
    ready: ComputedRef<boolean>
    /**
     * Computed indicating if the user is logged in.
     */
    loggedIn: ComputedRef<boolean>
    /**
     * The user object if logged in, null otherwise.
     */
    user: ComputedRef<PbUser | null>
    /**
     * The session object.
     */
    session: Ref<UserSession>
    /**
     * Fetch the user session from the server.
     */
    fetch: () => Promise<void>
    /**
     * Clear the user session and remove the session cookie.
     */
    clear: () => Promise<void>
    /**
     * Open the OAuth route in a popup that auto-closes when successful.
     */
    openInPopup: (route: string, size?: { width?: number, height?: number }) => void
  }

}
