import Symbol from "../Symbol/Symbol";
import { Service } from "../../models/service";

type Props = {
    service?: Service;
};

export default function ServiceLogo(props: Props) {
    const { service } = props;

    // @ts-ignore
    const iconURL = new URL(window.apiURL);
    iconURL.pathname = `/api/services/icons/${service?.icon}`;

    if (service?.icon) {
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

    return <Symbol name="extension" style={{ opacity: 0.8 }} />;
}
