import Popup from "../../../../components/Popup/Popup";
import { Button, Input, MaterialIcon } from "@vertex-center/components";
import { ChangeEvent, Fragment, useState } from "react";
import { APIError } from "../../../../components/Error/APIError";
import { useCreateContainer } from "../../hooks/useCreateContainer";
import { ProgressOverlay } from "../../../../components/Progress/Progress";

type Props = {
    show: boolean;
    dismiss: () => void;
};

export default function ManualInstallPopup(props: Readonly<Props>) {
    const { show, dismiss } = props;

    const [image, setImage] = useState<string>();

    const { createContainer, isCreatingContainer, errorCreatingContainer } =
        useCreateContainer({
            onSuccess: dismiss,
        });

    const create = () => createContainer({ image });

    const onImageChange = (e: ChangeEvent<HTMLInputElement>) => {
        setImage(e.target.value);
    };

    const actions = (
        <Fragment>
            <Button variant="outlined" onClick={dismiss}>
                Cancel
            </Button>
            <Button
                variant="colored"
                onClick={create}
                rightIcon={<MaterialIcon icon="download" />}
            >
                Install
            </Button>
        </Fragment>
    );

    return (
        <Popup
            show={show}
            onDismiss={dismiss}
            title="Install from Docker Registry"
            actions={actions}
        >
            <Input
                id="image"
                label="Image"
                placeholder="postgres"
                value={image}
                onChange={onImageChange}
                disabled={isCreatingContainer}
                required
            />
            <ProgressOverlay show={isCreatingContainer} />
            <APIError error={errorCreatingContainer} />
        </Popup>
    );
}
