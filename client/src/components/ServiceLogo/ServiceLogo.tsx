import { PuzzlePiece } from "@phosphor-icons/react";

type Props = {
    icon?: string;
    color?: string;
};

export default function ServiceLogo(props: Readonly<Props>) {
    const { icon, color } = props;

    // @ts-ignore
    const iconURL = new URL(window.api_urls.containers);
    iconURL.pathname = `/api/templates/icons/${icon}`;

    if (!icon) {
        return <PuzzlePiece size={32} style={{ opacity: 0.8 }} />;
    }

    if (icon.endsWith(".svg")) {
        return (
            <span
                style={{
                    display: "inline-block",
                    maskImage: `url(${iconURL.href})`,
                    backgroundColor: color,
                    width: 32,
                    height: 32,
                }}
            />
        );
    }

    return (
        <img
            alt="Service icon"
            src={iconURL.href}
            style={{
                width: 32,
                height: 32,
            }}
        />
    );
}
