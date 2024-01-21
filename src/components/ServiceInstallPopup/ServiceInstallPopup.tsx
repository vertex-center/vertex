import { APIError } from "../Error/APIError";
import { Button, MaterialIcon, Paragraph } from "@vertex-center/components";
import Popup from "../Popup/Popup";
import { Template as ServiceModel } from "../../apps/Containers/backend/template";
import { Fragment, useState } from "react";

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

    const actions = (
        <Fragment>
            <Button onClick={dismiss}>Cancel</Button>
            <Button
                variant="colored"
                onClick={install}
                rightIcon={<MaterialIcon icon="add" />}
                disabled={error !== undefined}
            >
                Create container
            </Button>
        </Fragment>
    );

    return (
        <Popup
            show={show}
            onDismiss={dismiss}
            title={`Install ${service?.name}`}
            actions={actions}
        >
            <Paragraph>{service?.description}</Paragraph>
            <APIError style={{ margin: 0 }} error={error} />
        </Popup>
    );
}
