<script>
    import { StoreAgentCredentials } from "../wailsjs/go/main/App.js";

    let email = "";
    let pw = "";
    let resultsVisible = false
    let success = false;

    async function updateCredentials() {
        let result = await StoreAgentCredentials(email.toString(), pw.toString())
        if (result == "error") {
            success = false;
        } else {
            success = true;
        }
        resultsVisible = true
    }

</script>
<h2>
    Enter your email and Zendesk password
</h2>
    
<input  id="email" name="email" bind:value={email} placeholder="Email"><br>
<input type="password" id="pw" name="pw" bind:value={pw} placeholder="Zendesk password"><br><br>    
<button on:click={updateCredentials}>Update</button>

{#if resultsVisible}
    {#if success}
        <p>Successfully updated your credentials</p>
    {:else}
        <p>Something went wrong</p>
    {/if}
{/if}


    
<style>
    button {
        border-radius: .25rem;
        cursor: pointer;
        font-weight: bold;
        background-color: #fff;
        box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
        max-width: 250px;
    }
    button:active {
        background-color: #efefef;
        box-shadow: 0 0px 0px 0 rgba(0, 0, 0, 0.1), 0 0px 0px 0 rgba(0, 0, 0, 0.06);	
    }
</style>