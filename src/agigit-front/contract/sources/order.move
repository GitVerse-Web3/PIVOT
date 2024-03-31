module nag::order {
    use std::signer;
    use std::vector;
    use aptos_framework::account;
    use aptos_framework::coin::{Self, Coin};
    use aptos_framework::event;
    use aptos_framework::aptos_coin::{AptosCoin};


    struct ModuleData has key {
        resource_signer_cap: account::SignerCapability,
    }

    struct Counter has key { i: u64 }

    struct Order has store {
        order_id: u64,
        status: u8,
        model_hash: u64,
        relay_pubkey_list: vector<u64>,
        down: u64,
        cost: u64,
        tip: u64,
        taker: address,
        product: u64,
    }

    struct OrderList has key {
        order_list: vector<Order>,
    }

    #[event]
    struct Pull has drop, store {
        address: address,
        order_id: u64,
        model_hash: u64,
        relay_pubkey_list: vector<u64>,
        down: u64,
        cost: u64,
        tip: u64,
    }

    #[event]
    struct RaisePayment has drop, store {
        address: address,
        order_id: u64,
        model_hash: u64,
        relay_pubkey_list: vector<u64>,
        down: u64,
        cost: u64,
        tip: u64,
    }

    #[event]
    struct TakeOrder has drop, store {
        relay_address: address,
        order_id: u64,
        order_address: address,
    }

    #[event]
    struct CompleteOrder has drop, store {
        order_address: address,
        order_id: u64,
        prime1: u64,
        prime2: u64,
    }


    fun init_module (account: &signer) {
        // store the capabilities within `ModuleData`
        let (signer, resource_signer_cap) = account::create_resource_account(account, vector[1,1,1,1,1]);
        move_to(account, ModuleData {
            resource_signer_cap,
        });
        coin::register<AptosCoin>(&signer);
        coin::register<AptosCoin>(account);
    }


    public entry fun pull(
        account: signer,
        model_hash: u64,
        relay_pubkey_list: vector<u64>,
        down: u64,
        cost: u64,
        tip: u64
    ) acquires OrderList, ModuleData {

        let amount1 = down;
        let amount2 = cost;
        let amount3 = tip;
        let sum = down + cost + tip;

        let module_data = borrow_global_mut<ModuleData>(@nag);
        let resource_signer = account::create_signer_with_capability(&module_data.resource_signer_cap);
        coin::register<AptosCoin>(&resource_signer);
        coin::register<AptosCoin>(&account);
        coin::transfer<AptosCoin>(&account, signer::address_of(&resource_signer), sum);

        let order_id = 0;
        let order = Order{
            order_id: order_id,
            status: 0,
            model_hash: model_hash,
            relay_pubkey_list: relay_pubkey_list,
            down: amount1,
            cost: amount2,
            tip: amount3,
            taker: signer::address_of(&account),
            product: 0,
        };
        if (exists<OrderList>(signer::address_of(&account))) {
            let order_list = &borrow_global<OrderList>(signer::address_of(&account)).order_list;
            order_id = vector::length<Order>(order_list);
            let order_list_mut = &mut borrow_global_mut<OrderList>(signer::address_of(&account)).order_list;
            vector::insert(order_list_mut, order_id, order);
        } else {
            let order_list: vector<Order> = vector[];
            vector::insert(&mut order_list, order_id, order);
            move_to(&account, OrderList{order_list});
        };


        event::emit(Pull{
            address: signer::address_of(&account),
            order_id: order_id,
            model_hash: model_hash,
            relay_pubkey_list: relay_pubkey_list,
            down: amount1,
            cost: amount2,
            tip: amount3,
        });
    }

    public entry fun raise_payment(
        account: signer,
        order_id: u64,
        down: u64,
        cost: u64,
        tip: u64
    ) acquires OrderList, ModuleData {
        assert!(exists<OrderList>(signer::address_of(&account)), 0);
        let order_list = &borrow_global<OrderList>(signer::address_of(&account)).order_list;
        let length = vector::length<Order>(order_list);
        assert!(length >= order_id, 1);

        let order_list_mut = &mut borrow_global_mut<OrderList>(signer::address_of(&account)).order_list;
        let order = vector::borrow_mut<Order>(order_list_mut, order_id);
        assert!(order.status == 0, 2);

        let amount1 = down;
        let amount2 = cost;
        let amount3 = tip;

        let sum = amount1 + amount2 + amount3;
        let module_data = borrow_global_mut<ModuleData>(@nag);
        let resource_signer = account::create_signer_with_capability(&module_data.resource_signer_cap);
        coin::transfer<AptosCoin>(&account, signer::address_of(&resource_signer), sum);

        order.down = order.down + amount1;
        order.cost = order.cost + amount2;
        order.tip = order.tip + amount3;

        let model_hash = order.model_hash;
        let relay_pubkey_list = order.relay_pubkey_list;

        event::emit(RaisePayment{
            address: signer::address_of(&account),
            order_id: order_id,
            model_hash: model_hash,
            relay_pubkey_list: relay_pubkey_list,
            down: amount1,
            cost: amount2,
            tip: amount3,
        });
    }

    public entry fun take_order(
        account: signer,
        order_id: u64,
        order_address: address,
        relay_pubkey: u64,
        x: u64,
    )  acquires OrderList, ModuleData {
        assert!(exists<OrderList>(order_address), 0);

        let order_list_mut = &mut borrow_global_mut<OrderList>(order_address).order_list;
        let order = vector::borrow_mut<Order>(order_list_mut, order_id);
        assert!(order.status == 0, 2);
        assert!(vector::contains(&order.relay_pubkey_list, &relay_pubkey), 3);
        order.status = 1;
        order.taker = signer::address_of(&account);
        order.product = x;

        let module_data = borrow_global_mut<ModuleData>(@nag);
        let resource_signer = account::create_signer_with_capability(&module_data.resource_signer_cap);
        let down_amount = order.down;

        event::emit(TakeOrder{
            relay_address: order.taker,
            order_id: order.order_id,
            order_address: order_address,
        });

        coin::transfer<AptosCoin>(&resource_signer,signer::address_of(&account), down_amount);

    }

    #[view]
    public fun get_order_status(addr: address, order_id: u64): u8 acquires OrderList {
        assert!(exists<OrderList>(addr), 0);
        let order_list = &borrow_global<OrderList>(addr).order_list;
        let length = vector::length<Order>(order_list);
        assert!(length > order_id, 1);
        let order = vector::borrow<Order>(order_list, order_id);
        order.status
    }

    public entry fun confirm_received(account: signer, order_id: u64) acquires OrderList {
        assert!(exists<OrderList>(signer::address_of(&account)), 0);
        let order_list = &borrow_global<OrderList>(signer::address_of(&account)).order_list;
        let length = vector::length<Order>(order_list);
        assert!(length > order_id, 1);

        let order_list_mut = &mut borrow_global_mut<OrderList>(signer::address_of(&account)).order_list;
        let order = vector::borrow_mut<Order>(order_list_mut, order_id);
        assert!(order.status == 1, 2);

        order.status = 2;
    }

    public entry fun complete_order(
        account: signer,
        order_id: u64,
        order_address: address,
        prime1: u64,
        prime2: u64
    ) acquires OrderList, ModuleData {
        assert!(exists<OrderList>(order_address), 0);
        let order_list = &borrow_global<OrderList>(order_address).order_list;
        let length = vector::length<Order>(order_list);
        assert!(length > order_id, 1);

        let order_list_mut = &mut borrow_global_mut<OrderList>(order_address).order_list;
        let order = vector::borrow_mut<Order>(order_list_mut, order_id);
        assert!(order.status == 2, 2);

        let x = prime1 * prime2;

        if (x == order.product) {
            order.status = 3;
            event::emit(CompleteOrder {
                order_address: order_address,
                order_id: order_id,
                prime1: prime1,
                prime2: prime2,
            });

            let module_data = borrow_global_mut<ModuleData>(@nag);
            let resource_signer = account::create_signer_with_capability(&module_data.resource_signer_cap);
            let amount = order.cost + order.tip;
            coin::transfer<AptosCoin>(&resource_signer,signer::address_of(&account), amount);
        } else {
            order.status = 4;
        }
    }
}
