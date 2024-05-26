import * as React from "react"
import {useState} from "react";

export interface AuthContextI {
    username: string
}

export const AuthContext = React.createContext<AuthContextI>(null)

const AuthContextProvider: React.FC<React.PropsWithChildren> = ({children}) => {
    const [authContext, setAuthContext] = useState()
    return <AuthContext.Provider value={authContext}>
        {children}
    </AuthContext.Provider>
}

export default AuthContextProvider