import Icon from "../Icon/Icon";
import { Service } from "../../models/service";

type Props = {
    service?: Service;
};

export default function ServiceLogo(props: Readonly<Props>) {
    const { service } = props;

    // @ts-ignore
    const iconURL = new URL(window.apiURL);
    iconURL.pathname = `/api/services/icons/${service?.icon}`;

    if (!service?.icon) {
        return <Icon name="extension" style={{ opacity: 0.8 }} />;
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
