import { Template } from "../../apps/Containers/backend/template";
import { MaterialIcon } from "@vertex-center/components";

type Props = {
    template?: Template;
};

export default function ServiceLogo(props: Readonly<Props>) {
    const { template } = props;

    // @ts-ignore
    const iconURL = new URL(window.api_urls.containers);
    iconURL.pathname = `/api/templates/icons/${template?.icon}`;

    if (!template?.icon) {
        return <MaterialIcon icon="extension" style={{ opacity: 0.8 }} />;
    }

    if (template?.icon.endsWith(".svg")) {
        return (
            <span
                style={{
                    maskImage: `url(${iconURL.href})`,
                    backgroundColor: template?.color,
                    width: 24,
                    height: 24,
                }}
            />
        );
    }

    return (
        <img
            alt="Service icon"
            src={iconURL.href}
            style={{
                width: 24,
                height: 24,
            }}
        />
    );
}
