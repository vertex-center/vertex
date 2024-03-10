import Popup, { PopupActions } from "../../../../components/Popup/Popup";
import { Container } from "../../backend/models";
import { Button, Horizontal } from "@vertex-center/components";
import styles from "../../../../components/TemplateInstallPopup/TemplateInstallPopup.module.sass";
import Spacer from "../../../../components/Spacer/Spacer";
import { CaretRight, CheckCircle } from "@phosphor-icons/react";
import { useNavigate } from "react-router-dom";
import TemplateImage from "../TemplateImage/TemplateImage";

type Props = {
    container: Container;
    onDismiss: () => void;
};

export default function ContainerInstalledPopup(props: Readonly<Props>) {
    const { container, onDismiss } = props;
    const navigate = useNavigate();

    if (!container) return null;

    const open = () => navigate(`/containers/${container.id}/logs`);

    const image = (
        <TemplateImage icon={container.icon} color={container.color} />
    );

    return (
        <Popup title={container?.name} onDismiss={onDismiss} image={image}>
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
                <Button onClick={onDismiss} variant="outlined">
                    Close
                </Button>
                <Button
                    variant="colored"
                    onClick={open}
                    rightIcon={<CaretRight />}
                >
                    Open
                </Button>
            </PopupActions>
        </Popup>
    );
}
