<script>
    //import { SearchContactByName, SearchContactByOrg } from "../wailsjs/go/main/App.js";
    //import Select from "./Select.svelte";
    import SiteSelect from "./OrgSelector.svelte";
    import CallerSelect from "./CallerSelector.svelte";
    import ProductSelector from "./ProductSelector.svelte";
    import { CreateTicket } from "../wailsjs/go/main/App";

    let callerName = "";
    let orgName = "";
    let subject = ""; //issue
    let problem = "";
    let assetName = "";
    let product = "";
    let phoneNumber = "";
    let issueStartDate = "";
    let rerender=true;
    let hasBeenClicked = false;

    async function submitNewTicket() {
        hasBeenClicked = true;
        if (isValidOrg && isValidCaller && isValidSubject && isValidProblem) {
            // create a new zendesk ticket
            let ticketId = await CreateTicket(callerName.toString(), "", phoneNumber.toString(), orgName.toString(), problem.toString(), subject.toString(), assetName.toString(), product.toString())
            if (ticketId == "Error creating ticket") {
                alert("There was an error\nno new ticket has been created")
            } else {
                alert("New ticket:\n" + ticketId)
            }
            
            // reset
            hasBeenClicked = false;

            // clear all fields
            orgName = ""
            callerName = ""
            subject = ""; //issue
            problem = "";
            assetName = "";
            phoneNumber = "";
            issueStartDate = "";
            rerender = !rerender
        }
        

    }

    $: isValidOrg = orgName != 'Organisation...';
    $: isValidCaller = callerName != 'Caller...';
    $: isValidSubject = subject.length > 5;
    $: isValidProblem = problem.length > 5;

</script>
<h2>
    Create a Zendesk ticket
</h2>
{#key rerender}
<table align="center">
    <tr>
        <td align="right"><label for="org">Organisation*</label></td>
        <td align="left"><SiteSelect bind:value={orgName} />{#if hasBeenClicked && !isValidOrg} 
            <p class="validation-error">Please select an organisation</p>
          {/if}</td>
    </tr>
    <tr>
        <td align="right"><label for="caller">Caller*</label></td>
        <td align="left"><CallerSelect bind:value={callerName}/>{#if hasBeenClicked && !isValidCaller} 
            <p class="validation-error">Please select a caller</p>
          {/if}</td>
    </tr>
    <tr>
        <td align="right"><label for="subj">Subject*</label></td>
        <td align="left"><input id="subj" bind:value={subject} placeholder="Issue title">{#if hasBeenClicked && !isValidSubject} 
            <p class="validation-error">Please enter a subject</p>
          {/if}</td>
    </tr>
    <tr>
        <td align="right"><label for="prob">Problem*</label></td>
        <td align="left"><textarea id="prob" bind:value={problem} placeholder="Detailed description of the problem"/>{#if hasBeenClicked && !isValidProblem} 
            <p class="validation-error">Please enter a problem description</p>
          {/if}</td>
    </tr>
    <tr>
        <td align="right"><label for="asset">Asset</label></td>
        <td align="left"><input id="asset" bind:value={assetName} placeholder="Name of the affected asset"/></td>
    </tr>
    <tr>
        <td align="right"><label for="prod">Product</label></td>
        <td align="left"><ProductSelector display_func={o => o.text} bind:value={product}/></td>
    </tr>
    <tr>
        <td align="right"><label for="phone">Caller's phone #</label></td>
        <td align="left"><input id="phone" bind:value={phoneNumber} placeholder="Caller's phone #"/></td>
    </tr>
    <tr>
        <td align="right"><label for="startDate">Date issue began</label></td>
        <td align="left"><input id="startDate" bind:value={issueStartDate} placeholder="Date the issue began"/></td>
    </tr>
</table>

{/key}
    
<button on:click={submitNewTicket}>
    Submit
</button>   

    
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
    textarea {
		width: 100%;
		height: 150px;
	}
    .validation-error {
    color: red;
    margin-top: 5px;
  }
</style>