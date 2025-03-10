import { invoke } from '@tauri-apps/api/core';

export const preventDefault = <T extends Event>(fn: (e: T) => void): ((e: T) => void) => {
	return (e: T) => {
		e.preventDefault();
		fn(e);
	};
};

export class GlobalState {
	private _state = $state({ name: '', greet: '' });
	private _state2 = $state({ inputMode: 'text' });

	get greet() {
		return this._state.greet;
	}
	set greet(value: string) {
		this._state.greet = value;
	}
	get name() {
		return this._state.name;
	}
	set name(value: string) {
		this._state.name = value;
	}
	get nlen() {
		return this.name.length;
	}

	get inputMode() {
		return this._state2.inputMode;
	}
	set inputMode(value: string) {
		this._state2.inputMode = value;
	}

	async submit() {
		this.greet = await invoke('greet', { name: this.name });
	}

	reset() {
		this.name = '';
		this.greet = '';
	}
}
