import { Service } from "../../models/service";
import { MaterialIcon } from "@vertex-center/components";

type Props = {
    service?: Service;
};

export default function ServiceLogo(props: Readonly<Props>) {
    const { service } = props;

    // @ts-ignore
    const iconURL = new URL(window.apiURL + ":7504");
    iconURL.pathname = `/api/services/icons/${service?.icon}`;

    if (!service?.icon) {
        return <MaterialIcon icon="extension" style={{ opacity: 0.8 }} />;
    }

    if (service?.icon.endsWith(".svg")) {
        return (
            <span
                style={{
                    maskImage: `url(${iconURL.href})`,
                    backgroundColor: service?.color,
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
