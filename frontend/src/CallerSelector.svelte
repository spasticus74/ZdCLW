<script>
    import { onMount } from "svelte";
    import { GetAllCallers } from "../wailsjs/go/main/App";

    export let callerOptions = [];
    export let display_func = a => a;
    export let index = 0;
    export let value;

    onMount(async () => {
        let callers = await GetAllCallers()
        let so = ['Caller...'];
        callers.forEach(caller => {
            so.push(caller)
        });
        callerOptions = so
    })
  
    $: {
      value = callerOptions[index];
      console.log(value);
    }
  </script>
  
  <select bind:value={index}>
    {#each callerOptions as callerOption, i}
      <option value={i}>{display_func(callerOption)}</option>
    {/each}
  </select>
  