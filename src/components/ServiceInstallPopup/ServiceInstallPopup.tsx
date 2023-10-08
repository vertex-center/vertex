import { Horizontal } from "../Layouts/Layouts";
import ServiceLogo from "../ServiceLogo/ServiceLogo";
import { Text, Title } from "../Text/Text";
import { APIError } from "../Error/APIError";
import Spacer from "../Spacer/Spacer";
import Button from "../Button/Button";
import Popup from "../Popup/Popup";
import { Service as ServiceModel } from "../../models/service";
import { useState } from "react";

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
        <Popup show={show} onDismiss={dismiss}>
            <Horizontal gap={15} alignItems="center">
                {service && <ServiceLogo service={service} />}
                <Title>{service?.name}</Title>
            </Horizontal>
            <Text>{service?.description}</Text>
            <APIError style={{ margin: 0 }} error={error} />
            <Horizontal gap={8}>
                <Spacer />
                <Button onClick={dismiss}>Cancel</Button>
                <Button
                    onClick={install}
                    primary
                    rightIcon="add"
                    disabled={error !== undefined}
                >
                    Create instance
                </Button>
            </Horizontal>
        </Popup>
    );
}
