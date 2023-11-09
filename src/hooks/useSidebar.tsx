import { ReactNode } from "react";
import { createPortal } from "react-dom";
import { useHasSidebar } from "@vertex-center/components";

export const useSidebar = (sidebar: ReactNode) => {
    useHasSidebar(true);

    const sidebarContainer = document.getElementsByClassName("app-sidebar")[0];
    let s = null;
    if (sidebarContainer) {
        s = createPortal(sidebar, sidebarContainer);
    }
    return s;
};
