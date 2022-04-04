import { TW } from "@utils/tailwind-mixin";
import type { TemplateResult } from "lit";
import { LitElement, html } from "lit";
import { customElement, property } from "lit/decorators.js";

@customElement("button-counter")
export class Counter extends TW(LitElement) {
  @property({ type: Number })
  private counter = 0;

  render(): TemplateResult {
    return html`
      <p class="text-red-500">Hello, ${this.counter}</p>
      <button class="block m-5" @click="${() => (this.counter += 1)}">Click</button>
    `;
  }
}
