using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using Zenject;

public class ClickManager : MonoBehaviour
{
	[Inject]
	BranchManager _branchManager;

	[Inject]
	public new Camera camera;


	RaycastHit hit;

	void rebase(Node node)
	{
		node.rebaseToMaster(_branchManager.masterHead);
	}



	void Start()
	{


	}

	// Update is called once per frame
	void Update()
	{
		Ray ray = camera.ScreenPointToRay(Input.mousePosition);

		if (Physics.Raycast(ray, out hit))
		{
			Transform objectHit = hit.transform;

			Node node = objectHit.GetComponent<Node>();

			rebase(node);
		}
	}
}
