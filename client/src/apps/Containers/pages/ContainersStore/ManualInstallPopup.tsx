import Popup from "../../../../components/Popup/Popup";
import { Button, FormItem, Input } from "@vertex-center/components";
import { ChangeEvent, Fragment, useState } from "react";
import { APIError } from "../../../../components/Error/APIError";
import { useCreateContainer } from "../../hooks/useCreateContainer";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { DownloadSimple } from "@phosphor-icons/react";

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
                rightIcon={<DownloadSimple />}
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
            <FormItem label="Image" required>
                <Input
                    placeholder="postgres"
                    value={image}
                    onChange={onImageChange}
                    disabled={isCreatingContainer}
                />
            </FormItem>
            <ProgressOverlay show={isCreatingContainer} />
            <APIError error={errorCreatingContainer} />
        </Popup>
    );
}
