import { APIError } from "../Error/APIError";
import { Button, Paragraph } from "@vertex-center/components";
import Popup, { PopupActions } from "../Popup/Popup";
import { Template as ServiceModel } from "../../apps/Containers/backend/template";
import { useState } from "react";
import { Plus } from "@phosphor-icons/react";
import ServiceLogo from "../ServiceLogo/ServiceLogo";
import styles from "./TemplateInstallPopup.module.sass";

type Props = {
    service: ServiceModel;
    show: boolean;
    dismiss: () => void;
    install: () => void;
};

export default function TemplateInstallPopup(props: Readonly<Props>) {
    const { service, show } = props;

    const [error, setError] = useState();

    const dismiss = () => {
        setError(undefined);
        props.dismiss();
    };

    const install = () => {
        setError(undefined);
        props.install();
    };

    const backgroundColor = `${service?.color ?? "#000000"}08`;

    return (
        <Popup
            show={show}
            onDismiss={dismiss}
            title={service?.name}
            image={
                <div className={styles.logo} style={{ backgroundColor }}>
                    <ServiceLogo template={service} />
                </div>
            }
        >
            <Paragraph>{service?.description}</Paragraph>
            <APIError style={{ margin: 0 }} error={error} />
            <PopupActions>
                <Button onClick={dismiss}>Cancel</Button>
                <Button
                    variant="colored"
                    onClick={install}
                    rightIcon={<Plus />}
                    disabled={error !== undefined}
                >
                    Create container
                </Button>
            </PopupActions>
        </Popup>
    );
}
