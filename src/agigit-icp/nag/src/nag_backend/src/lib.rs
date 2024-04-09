use std::borrow::BorrowMut;
use std::cell::RefCell;
use candid::{types::number::Nat, Principal};
use ic_cdk::{caller,query, update};
use ic_ledger_types::{
    AccountIdentifier, BlockIndex, Memo, Subaccount, Tokens, DEFAULT_SUBACCOUNT,
    MAINNET_LEDGER_CANISTER_ID,
};
use std::collections::BTreeMap;


// struct Order {
//     order_id: u64,
//     status: u8,
//     model_hash: u64,
//     relay_pubkey_list: Vec<u64>,
//     down: u64,
//     cost: u64,
//     tip: u64,
//     taker: Principal,
//     product: u64,
// }

struct Order {
    order_id: Nat,
    model_hash: Nat,
    status: u8,
}

type OrderStore = BTreeMap<u64, Order>;

thread_local! {
    //static COUNTER: RefCell<Nat> = RefCell::new(Nat::from(0));
    static ORDERSTORE: RefCell<OrderStore> = RefCell::default();
}

#[update]
fn pull(model_hash: Nat) {
    let order = Order {
        order_id: Nat::from(1u8),
        model_hash,
        status: 0,
    };
    ORDERSTORE.with(|order_store| {
        order_store.borrow_mut().insert(1, order)
    });

}

#[update]
fn raise_payment(order_id: Nat) {

}

#[update]
fn take_order(order_id: Nat) {
    ORDERSTORE.with(|order_store| {
        order_store
            .borrow_mut()
            .get_mut(&1)
            .and_then(|order| {
                order.status = 1;
                Some(())
            });

    });
}

#[update]
fn confirm_received(order_id: Nat) {
    ORDERSTORE.with(|order_store| {
        order_store
            .borrow_mut()
            .get_mut(&1)
            .and_then(|order| {
                order.status = 2;
                Some(())
            });

    });
}

#[update]
fn complete_order(order_id: Nat) {
    ORDERSTORE.with(|order_store| {
        order_store
            .borrow_mut()
            .get_mut(&1)
            .and_then(|order| {
                order.status = 3;
                Some(())
            });

    });
}

#[query]
fn get_order_status(order_id: Nat) -> Nat {
    Nat::from(ORDERSTORE.with(|order_store| {
        order_store.borrow().get(&1).unwrap().status
    }))
}





// async fn deposit(caller: Principal, amount: Tokens) -> Result<Nat, DepositErr> {
//     let canister_id = ic_cdk::api::id();
//
//     let transfer_args = ic_ledger_types::TransferArgs {
//         memo: Memo(0),
//         amount: amount,
//         fee: Tokens::from_e8s(10000),
//         from_subaccount: Some(principal_to_subaccount(&caller)),
//         to: AccountIdentifier::new(&canister_id, &DEFAULT_SUBACCOUNT),
//         created_at_time: None,
//     };
//     ic_ledger_types::transfer(ledger_canister_id, transfer_args)
//         .await
//         .map_err(|_| DepositErr::TransferFailure)?
//         .map_err(|_| DepositErr::TransferFailure)?;
//
//     ic_cdk::println!(
//         "Deposit of {} ICP in account {:?}",
//         balance - Tokens::from_e8s(ICP_FEE),
//         &account
//     );
//
//     Ok((balance.e8s() - ICP_FEE).into())
// }