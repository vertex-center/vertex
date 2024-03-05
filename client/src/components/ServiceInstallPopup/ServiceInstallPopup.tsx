import { APIError } from "../Error/APIError";
import { Button, Paragraph } from "@vertex-center/components";
import Popup, { PopupActions } from "../Popup/Popup";
import { Template as ServiceModel } from "../../apps/Containers/backend/template";
import { Fragment, useState } from "react";
import { Plus } from "@phosphor-icons/react";

type Props = {
    service: ServiceModel;
    show: boolean;
    dismiss: () => void;
    install: () => void;
};

export default function ServiceInstallPopup(props: Readonly<Props>) {
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

    return (
        <Popup
            show={show}
            onDismiss={dismiss}
            title={`Install ${service?.name}`}
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
