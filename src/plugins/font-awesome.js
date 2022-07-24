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
  faArrowDown,
  faArrowUp,
  faCheckCircle,
  faFileDownload,
  faFileInvoice,
  faFileSignature,
  faFileUpload,
  faLink,
  faMailBulk,
  faNewspaper,
  faPaperPlane,
  faPaste,
  faPenNib,
  faPlus,
  faSlidersH,
  faSyncAlt,
  faTrashAlt,
  faUserCircle,
  faUserClock,
} from "@fortawesome/free-solid-svg-icons";

library.add(
  faArrowDown,
  faArrowUp,
  faCheckCircle,
  faCopy,
  faFileCode,
  faFileDownload,
  faFileInvoice,
  faFileSignature,
  faFileUpload,
  faFileWord,
  faLink,
  faMailBulk,
  faNewspaper,
  faPaperPlane,
  faPaste,
  faPenNib,
  faPlus,
  farNewspaper,
  faSlidersH,
  faSyncAlt,
  faTrashAlt,
  faUserCircle,
  faUserClock
);

Vue.component("FontAwesomeIcon", FontAwesomeIcon);
