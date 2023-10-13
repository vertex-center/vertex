import { Service as ServiceModel } from "../../models/service";

import styles from "./Service.module.sass";
import { Caption } from "../Text/Text";
import Progress from "../Progress";
import ServiceLogo from "../ServiceLogo/ServiceLogo";

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
        <div className={styles.service} onClick={onInstall}>
            <div className={styles.logo}>
                <ServiceLogo service={service} />
            </div>
            <div className={styles.info}>
                <div>{service?.name}</div>
                <Caption>{service?.description}</Caption>
            </div>
            {installedCountText && (
                <Caption className={styles.count}>{installedCountText}</Caption>
            )}
            {downloading && <Progress infinite />}
        </div>
    );
}
