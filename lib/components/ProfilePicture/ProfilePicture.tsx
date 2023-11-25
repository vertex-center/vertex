import "./ProfilePicture.sass";
import cx from "classnames";
import { HTMLProps } from "react";
import { MaterialIcon } from "../MaterialIcon/MaterialIcon";

export type ProfilePictureProps = HTMLProps<HTMLImageElement> & {
    size?: number;
};

export function ProfilePicture(props: Readonly<ProfilePictureProps>) {
    const { className, alt, size = 40, ...others } = props;

    const properties = {
        className: cx("profile-picture", className),
        ...others,
    };

    if (props.src === undefined) {
        return (
            <div style={{ width: size, height: size }} {...properties}>
                <MaterialIcon icon="person" />
            </div>
        );
    }
    return <img alt={alt} width={size} height={size} {...properties} />;
}
