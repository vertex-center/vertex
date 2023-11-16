import { Horizontal } from "../Layouts/Layouts";
import ServiceLogo from "../ServiceLogo/ServiceLogo";
import { Title } from "../Text/Text";
import { APIError } from "../Error/APIError";
import Spacer from "../Spacer/Spacer";
import { Button, MaterialIcon, Paragraph } from "@vertex-center/components";
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
            <Paragraph>{service?.description}</Paragraph>
            <APIError style={{ margin: 0 }} error={error} />
            <Horizontal gap={8}>
                <Spacer />
                <Button onClick={dismiss}>Cancel</Button>
                <Button
                    variant="colored"
                    onClick={install}
                    rightIcon={<MaterialIcon icon="add" />}
                    disabled={error !== undefined}
                >
                    Create container
                </Button>
            </Horizontal>
        </Popup>
    );
}
