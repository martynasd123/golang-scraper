import * as React from "react"
import {PropsWithChildren, useEffect, useState} from "react";
import {useLocation} from "react-router-dom";

interface HeaderContext {
    setProgressBarProgress: (progress: number | null) => void;
    setBackButtonLink: (url: string | null) => void;
    backButtonLink: string | null;
    progressBarProgress: number | null;
}

export const HeaderContext = React.createContext<HeaderContext>(null);

export const HeaderContextProvider: React.FC<PropsWithChildren> = ({children}) => {
    const [progressBarProgress, setProgressBarProgress] = useState(null)
    const [backButtonLink, setBackButtonLink] = useState(null)

    const location = useLocation();

    useEffect(() => {
        setProgressBarProgress(null)
        setBackButtonLink(null)
    }, [location]);

    return <HeaderContext.Provider value={{
        setProgressBarProgress: (progress) => setProgressBarProgress(progress),
        setBackButtonLink: (link) => setBackButtonLink(link),
        backButtonLink,
        progressBarProgress,
    }}>
        {children}
    </HeaderContext.Provider>
}