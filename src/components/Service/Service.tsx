import { Service as ServiceModel } from "../../models/service";

import styles from "./Service.module.sass";
import { Caption } from "../Text/Text";
import Progress from "../Progress";
import ServiceLogo from "../ServiceLogo/ServiceLogo";
import ListItem from "../List/ListItem";
import ListIcon from "../List/ListIcon";
import ListInfo from "../List/ListInfo";
import ListTitle from "../List/ListTitle";
import ListActions from "../List/ListActions";
import Icon from "../Icon/Icon";

type Props = {
    service: ServiceModel;
    onInstall: () => void;
    downloading?: boolean;
    installedCount?: number;
};

export default function Service(props: Readonly<Props>) {
    const { service, onInstall, downloading, installedCount } = props;

    let installedCountText = "";
    if (installedCount === 1) {
        installedCountText = "Installed in 1 container";
    } else if (installedCount > 1) {
        installedCountText = `Installed in ${installedCount} containers`;
    }

    return (
        <ListItem onClick={onInstall}>
            <ListIcon>
                <ServiceLogo service={service} />
            </ListIcon>
            <ListInfo>
                <ListTitle>{service?.name}</ListTitle>
                <Caption>{service?.description}</Caption>
            </ListInfo>
            <ListActions>
                {installedCountText && (
                    <Caption className={styles.count}>
                        {installedCountText}
                    </Caption>
                )}
                {downloading && <Progress infinite />}
                <Icon name="download" />
            </ListActions>
        </ListItem>
    );
}
