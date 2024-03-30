using System.Collections;
using System.Collections.Generic;
using System.Linq;
using UnityEngine;
using Zenject;

public class _BranchManager : PrefabFactory<ICommit>, IFactory<ICommit>, BranchManager
{
	IFactory<string, ICommit, ICommit> _nodeFactory;
	DiContainer _container;
	ICommit _masterHead;

	public _BranchManager(IFactory<string, ICommit, ICommit> nodeFactory, DiContainer container)
	{
		_nodeFactory = nodeFactory;
		_container = container;

		this.masterHead = null;
		this.allNodes = new List<ICommit>();


	}

	public ICommit masterHead
	{
		get
		{
			return _masterHead;
		}
		set
		{

			_masterHead = value;
			if (_masterHead != null)
			{
				_masterHead.chosenToBeHead();

			}


		}
	}

	public IList<ICommit> allNodes { get; }

	public IList<ICommit> newHeadCandidates
	{
		get
		{
			var ans = new List<ICommit>();
			foreach (var item in allNodes)
			{
				if (item.parentModel == masterHead)
				{
					ans.Add(item);

				}

			}

			return ans;
		}
	}


	public ICommit Create()
	{
		ICommit x = _nodeFactory.Create("Fuck OpenAI", masterHead);


		Node prefab = Resources.Load<Node>("node");
		Node ans = (Node)Create(prefab);
		ans._commit = x;
		if (ans.parentModel == null)
		{
			ans.transform.position = Vector3.zero;
		}
		else
		{
			ans.transform.position = new Vector3(
				Random.Range(-ans.r * 2, ans.r * 2),
				Random.Range(-ans.r * 2, ans.r * 2),
				Random.Range(-ans.r * 2, ans.r * 2)
				);
			ans.updateY((Node)masterHead);

		}


		allNodes.Add(ans);
		return ans;
	}


	public ICommit createNewNode()
	{
		return Create();
	}

	public void Init()
	{
		ICommit root = createNewNode();
		masterHead = root;


	}

	public void newMasterHead()
	{
		ICommit ans = newHeadCandidates.First();
		foreach (var item in newHeadCandidates)
		{
			if (item.compressionRatio < ans.compressionRatio)
			{
				ans = item;
			}

		}

		masterHead = ans;
	}
}
