import { TW } from "@utils/tailwind-mixin";
import type { TemplateResult } from "lit";
import { html, LitElement } from "lit";
import { customElement } from "lit/decorators.js";
import { unsafeHTML } from 'lit/directives/unsafe-html.js';

import icon from "/icons/adjustments.svg?raw";

import "@components/button-counter";
import "@styles/main.css";

@customElement("index-page")
export class IndexPage extends TW(LitElement) {
  render(): TemplateResult {
    return html`
      <div class="container">
        ${unsafeHTML(icon)}
        <button-counter></button-counter>
      </div>
    `;
  }
}
