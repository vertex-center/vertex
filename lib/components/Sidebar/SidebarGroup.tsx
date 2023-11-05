import { PropsWithChildren } from "react";
import { Title } from "../Title/Title.tsx";

export type SidebarGroupProps = PropsWithChildren<{
    title?: string;
}>;

export function SidebarGroup(props: Readonly<SidebarGroupProps>) {
    const { title, children } = props;

    return (
        <div className="sidebar-group">
            {title && (
                <Title variant="h5" className="sidebar-title">
                    {title}
                </Title>
            )}
            {children}
        </div>
    );
}
