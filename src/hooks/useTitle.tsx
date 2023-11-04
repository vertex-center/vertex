import { useContext, useEffect } from "react";
import { HeaderContext } from "../components/Header/Header";

export const useTitle = (title: string) => {
    const { setTitle } = useContext(HeaderContext);

    useEffect(() => {
        setTitle(title);
        return () => setTitle(undefined);
    }, [title]);
};
