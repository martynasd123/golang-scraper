import * as React from "react"
import {Dispatch, SetStateAction, useEffect, useState} from "react";

export interface AuthContextI {
    username: string,
    setUsername: Dispatch<SetStateAction<string>>,
    isAuthenticated: boolean
}

export const AuthContext = React.createContext<AuthContextI>(null)

const AuthContextProvider: React.FC<React.PropsWithChildren> = ({children}) => {
    const [username, setUsername] = useState(localStorage.getItem("username"))

    useEffect(() => {
        if (username) {
            localStorage.setItem("username", username)
        } else {
            localStorage.removeItem("username")
        }
    }, [username])

    return <AuthContext.Provider value={{username, setUsername, isAuthenticated: !!username}}>
        {children}
    </AuthContext.Provider>
}

export default AuthContextProvider