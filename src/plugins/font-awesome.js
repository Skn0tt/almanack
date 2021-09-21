import Vue from "vue";
import { library } from "@fortawesome/fontawesome-svg-core";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import {
  faCopy,
  faFileCode,
  faFileWord,
  faNewspaper as farNewspaper,
} from "@fortawesome/free-regular-svg-icons";
import {
  faFileDownload,
  faLink,
  faFileUpload,
  faFileInvoice,
  faFileSignature,
  faPaperPlane,
  faMailBulk,
  faNewspaper,
  faPenNib,
  faSlidersH,
  faSyncAlt,
  faTrashAlt,
  faUserCircle,
  faUserClock,
} from "@fortawesome/free-solid-svg-icons";

// Buefy icons
import {
  faAngleDown,
  faAngleLeft,
  faAngleRight,
  faArrowUp,
  faCaretDown,
  faCaretUp,
  faCheck,
  faCheckCircle,
  faExclamationCircle,
  faExclamationTriangle,
  faEye,
  faEyeSlash,
  faFileImage,
  faInfoCircle,
  faMinus,
  faPlus,
  faTimesCircle,
  faUpload,
} from "@fortawesome/free-solid-svg-icons";

library.add(
  faAngleDown,
  faAngleLeft,
  faAngleRight,
  faArrowUp,
  faCaretDown,
  faCaretUp,
  faCheck,
  faCheckCircle,
  faCopy,
  faExclamationCircle,
  faExclamationTriangle,
  faEye,
  faEyeSlash,
  faFileCode,
  faFileDownload,
  faFileImage,
  faFileInvoice,
  faFileSignature,
  faFileUpload,
  faFileWord,
  faInfoCircle,
  faLink,
  faMailBulk,
  faMinus,
  faNewspaper,
  faPaperPlane,
  faPenNib,
  faPlus,
  farNewspaper,
  faSlidersH,
  faSyncAlt,
  faTimesCircle,
  faTrashAlt,
  faUpload,
  faUserCircle,
  faUserClock
);

Vue.component("font-awesome-icon", FontAwesomeIcon);
