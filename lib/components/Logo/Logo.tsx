import Icon from "../../assets/logo.svg";
import { ImgHTMLAttributes } from "react";

export type LogoProps = ImgHTMLAttributes<HTMLImageElement>;

export function Logo(props: Readonly<LogoProps>) {
    const { alt = "App Logo", ...others } = props;
    return <img width={30} src={Icon} alt={alt} {...others} />;
}
