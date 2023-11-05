import "./Header.sass";
import { HTMLProps, ReactNode, useContext } from "react";
import cx from "classnames";
import { Logo } from "../Logo/Logo.tsx";
import { Title } from "../Title/Title.tsx";
import { Link, LinkProps } from "../Link/Link.tsx";
import { PageContext } from "../../contexts/PageContext";

interface IHeaderLink {
    className?: string;
}

export type HeaderProps<T, U> = HTMLProps<HTMLDivElement> & {
    appName?: string;
    leading?: ReactNode;
    linkLogo?: LinkProps<T>;
    linkBack?: LinkProps<U>;
};

export function Header<T extends IHeaderLink, U extends IHeaderLink>(
    props: Readonly<HeaderProps<T, U>>,
) {
    const {
        className,
        linkLogo,
        linkBack,
        appName,
        leading,
        children,
        ...others
    } = props;

    const { className: classNameLinkLogo, ...linkLogoProps } = linkLogo ?? {};
    const { className: classNameLinkBack, ...linkBackProps } = linkBack ?? {};

    const { title } = useContext(PageContext);

    const leadingElement = leading && (
        <Link
            className={cx("header-leading-link", classNameLinkBack)}
            {...linkBackProps}
        >
            <div className="header-leading">{leading}</div>
        </Link>
    );

    return (
        <header className={cx("header", className)} {...others}>
            <div className="header-top">
                {leadingElement}
                <Link
                    className={cx("header-logo", classNameLinkLogo)}
                    {...linkLogoProps}
                >
                    <Logo size={24} />
                    <Title variant="h4">{appName ?? "Vertex"}</Title>
                </Link>
            </div>
            {(children || title) && (
                <div className="header-bottom">
                    {children}
                    {!children && title && <Title variant="h1">{title}</Title>}
                </div>
            )}
        </header>
    );
}
