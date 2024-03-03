import { Template as TemplateModel } from "../../apps/Containers/backend/template";

import styles from "./Service.module.sass";
import { Caption } from "../Text/Text";
import Progress from "../Progress";
import ServiceLogo from "../ServiceLogo/ServiceLogo";
import {
    ListActions,
    ListIcon,
    ListInfo,
    ListItem,
    ListTitle,
} from "@vertex-center/components";
import { DownloadSimple } from "@phosphor-icons/react";

type Props = {
    template: TemplateModel;
    onInstall: () => void;
    downloading?: boolean;
    installedCount?: number;
};

export default function Service(props: Readonly<Props>) {
    const { template, onInstall, downloading, installedCount } = props;

    let installedCountText = "";
    if (installedCount === 1) {
        installedCountText = "Installed in 1 container";
    } else if (installedCount > 1) {
        installedCountText = `Installed in ${installedCount} containers`;
    }

    return (
        <ListItem onClick={onInstall}>
            <ListIcon>
                <ServiceLogo template={template} />
            </ListIcon>
            <ListInfo>
                <ListTitle>{template?.name}</ListTitle>
                <Caption>{template?.description}</Caption>
            </ListInfo>
            <ListActions>
                {installedCountText && (
                    <Caption className={styles.count}>
                        {installedCountText}
                    </Caption>
                )}
                {downloading && <Progress infinite />}
                <DownloadSimple />
            </ListActions>
        </ListItem>
    );
}
