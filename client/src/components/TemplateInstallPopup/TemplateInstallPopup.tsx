import { APIError } from "../Error/APIError";
import { Button, Paragraph, Spinner } from "@vertex-center/components";
import Popup, { PopupActions } from "../Popup/Popup";
import { Template as ServiceModel } from "../../apps/Containers/backend/template";
import { useState } from "react";
import { Plus } from "@phosphor-icons/react";
import { useCreateContainer } from "../../apps/Containers/hooks/useCreateContainer";
import Spacer from "../Spacer/Spacer";
import { Container } from "../../apps/Containers/backend/models";
import ContainerInstalledPopup from "../../apps/Containers/components/ContainerInstalledPopup/ContainerInstalledPopup";
import TemplateImage from "../../apps/Containers/components/TemplateImage/TemplateImage";

type Props = {
    service: ServiceModel;
    dismiss: () => void;
};

export default function TemplateInstallPopup(props: Readonly<Props>) {
    const { service, dismiss: _dismiss } = props;

    const [container, setContainer] = useState<Container>(null);

    const dismiss = () => {
        _dismiss();
    };

    const { createContainer, isCreatingContainer, errorCreatingContainer } =
        useCreateContainer({
            onSuccess: (data, error, options) => {
                setContainer(data as Container);
            },
        });

    if (container) {
        return (
            <ContainerInstalledPopup
                container={container}
                onDismiss={dismiss}
            />
        );
    }

    const image = <TemplateImage color={service?.color} icon={service?.icon} />;

    return (
        <Popup onDismiss={dismiss} title={service?.name} image={image}>
            <Paragraph>{service?.description}</Paragraph>
            <APIError style={{ margin: 0 }} error={errorCreatingContainer} />
            <PopupActions>
                {isCreatingContainer && <Spinner label="Creating..." />}
                <Spacer />
                <Button onClick={dismiss} disabled={isCreatingContainer}>
                    Cancel
                </Button>
                <Button
                    variant="colored"
                    onClick={() =>
                        createContainer({
                            template_id: service.id,
                        })
                    }
                    rightIcon={<Plus />}
                    disabled={isCreatingContainer}
                >
                    Create container
                </Button>
            </PopupActions>
        </Popup>
    );
}
