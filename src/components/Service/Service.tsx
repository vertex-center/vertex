import { Service as ServiceModel } from "../../models/service";

import styles from "./Service.module.sass";
import { Caption } from "../Text/Text";
import { Horizontal, Vertical } from "../Layouts/Layouts";
import Symbol from "../Symbol/Symbol";
import Spacer from "../Spacer/Spacer";
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
        <Horizontal className={styles.service} gap={16} alignItems="center">
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
            <Vertical gap={6}>
                <div>{service?.name}</div>
                <Caption>{service?.description}</Caption>
            </Vertical>
            <Spacer />
            {installedCountText && <Caption>{installedCountText}</Caption>}
            {downloading && <Progress infinite />}
            {!downloading && (
                <Button rightSymbol="add" onClick={onInstall}>
                    New instance
                </Button>
            )}
        </Horizontal>
    );
}
