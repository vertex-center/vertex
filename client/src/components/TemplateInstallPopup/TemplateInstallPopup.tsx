import { APIError } from "../Error/APIError";
import {
    Button,
    Horizontal,
    Paragraph,
    Spinner,
} from "@vertex-center/components";
import Popup, { PopupActions } from "../Popup/Popup";
import { Template as ServiceModel } from "../../apps/Containers/backend/template";
import { useState } from "react";
import { CaretRight, CheckCircle, Plus } from "@phosphor-icons/react";
import ServiceLogo from "../ServiceLogo/ServiceLogo";
import styles from "./TemplateInstallPopup.module.sass";
import { useCreateContainer } from "../../apps/Containers/hooks/useCreateContainer";
import Spacer from "../Spacer/Spacer";
import { Container } from "../../apps/Containers/backend/models";
import { useNavigate } from "react-router-dom";

type Props = {
    service: ServiceModel;
    show: boolean;
    dismiss: () => void;
};

export default function TemplateInstallPopup(props: Readonly<Props>) {
    const { service, show } = props;

    const navigate = useNavigate();

    const [state, setState] = useState<"installing" | "installed">(
        "installing"
    );
    const [container, setContainer] = useState<Container>(null);

    const dismiss = () => {
        setState("installing");
        props.dismiss();
    };

    const open = () => {
        navigate(`/containers/${container.id}/logs`);
    };

    const { createContainer, isCreatingContainer, errorCreatingContainer } =
        useCreateContainer({
            onSuccess: (data, error, options) => {
                setState("installed");
                setContainer(data as Container);
            },
        });

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
            <APIError style={{ margin: 0 }} error={errorCreatingContainer} />
            {state === "installing" && (
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
            )}
            {state === "installed" && (
                <PopupActions>
                    <Horizontal
                        gap={8}
                        alignItems="center"
                        className={styles.installed}
                    >
                        <CheckCircle />
                        Installed successfully
                    </Horizontal>
                    <Spacer />
                    <Button
                        onClick={dismiss}
                        disabled={isCreatingContainer}
                        variant="outlined"
                    >
                        Close
                    </Button>
                    <Button
                        variant="colored"
                        onClick={open}
                        rightIcon={<CaretRight />}
                        disabled={isCreatingContainer}
                    >
                        Open
                    </Button>
                </PopupActions>
            )}
        </Popup>
    );
}
