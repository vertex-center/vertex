import { Service as ServiceModel } from "../../models/service";

import styles from "./Service.module.sass";
import { Caption } from "../Text/Text";
import Symbol from "../Symbol/Symbol";
import Button from "../Button/Button";
import Progress from "../Progress";

type Props = {
    service: ServiceModel;
    onInstall: () => void;
    downloading?: boolean;
    installedCount?: number;
};

export default function Service(props: Props) {
    const { service, onInstall, downloading, installedCount } = props;

    let installedCountText = "";
    if (installedCount === 1) {
        installedCountText = "Installed in 1 instance";
    } else if (installedCount > 1) {
        installedCountText = `Installed in ${installedCount} instances`;
    }

    // @ts-ignore
    const iconURL = new URL(window.apiURL);
    iconURL.pathname = `/api/services/icons/${service?.icon}`;

    return (
        <div className={styles.service}>
            <div className={styles.logo}>
                {service?.icon ? (
                    <span
                        style={{
                            maskImage: `url(${iconURL.href})`,
                            backgroundColor: service?.color,
                            width: 24,
                            height: 24,
                        }}
                    />
                ) : (
                    <Symbol name="extension" style={{ opacity: 0.8 }} />
                )}
            </div>
            <div className={styles.info}>
                <div>{service?.name}</div>
                <Caption>{service?.description}</Caption>
            </div>
            {installedCountText && (
                <Caption className={styles.count}>{installedCountText}</Caption>
            )}
            {downloading && <Progress infinite />}
            {!downloading && (
                <Button
                    className={styles.add}
                    rightSymbol="add"
                    onClick={onInstall}
                >
                    New instance
                </Button>
            )}
        </div>
    );
}
