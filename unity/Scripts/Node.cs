using System;
using System.Collections;
using System.Collections.Generic;
using TMPro;
using UnityEngine;
using Zenject;

public class Node : MonoBehaviour, ICommit, IInitializable
{


	public float scaleFactor = 10;
	public float speed = 0.1f;
	public float r = 6;
	public float deltaY = 1;

	public ICommit _commit;

	[SerializeField]
	TextMeshPro _textMesh;
	[SerializeField]
	LineRenderer _line;


	Renderer _renderer;


	public long modelHashID => _commit.modelHashID;

	public PublicKey author => _commit.author;

	public long authorSignature => _commit.authorSignature;

	public DateTime timestamp => _commit.timestamp;

	public string commitMessage => _commit.commitMessage;

	public ICommit parentModel
	{
		get => _commit.parentModel;
		set
		{

			_commit.parentModel = value;



		}
	}

	public double compressionRatio
	{
		get => _commit.compressionRatio;
		set
		{
			_commit.compressionRatio = value;
			updateScale();
		}
	}

	public new Tag tag => _commit.tag;

	public bool checkValid()
	{
		return _commit.checkValid();
	}

	public void chosenToBeHead()
	{
		_commit.chosenToBeHead();






		_renderer.material.color = Color.red;
		if (parentModel != null)
		{
			((Node)parentModel)._renderer.material.color = Color.blue;

		}

		updateLine();
	}

	public byte[] getFullModel()
	{
		return _commit.getFullModel();
	}

	public void rebaseToMaster(ICommit head)
	{

		if (!this.tag.isMaster)
		{
			Debug.Log("rebase");
			Node h = (Node)head;

			_commit.rebaseToMaster(head);


			updateScale();
			updateY(h);
			updateLine();
		}

	}

	public void updateLine()
	{
		if (parentModel != null)
		{

			float c = (float)this.compressionRatio;
			_line.startWidth = c;
			_line.endWidth = c;
			_line.SetPositions(new Vector3[] { this.transform.position, ((Node)parentModel).transform.position });
		}
		else
		{
			_line.enabled = false;
		}
	}


	float myY = 0;
	public void updateY(Node head)
	{
		var v = head.transform.position;
		float dy = head.transform.localScale.y / 2 + deltaY * (float)this.compressionRatio + this.transform.localScale.y / 2;
		v.y += dy;
		v.x = this.transform.position.x;
		v.z = this.transform.position.z;
		this.transform.position = v;
		myY = v.y;
	}

	// Start is called before the first frame update


	// Update is called once per frame
	void LateUpdate()
	{
		Vector3 target;
		var v = this.transform.position;
		if (!this.tag.isMaster)
		{
			float c = (float)this.parentModel.compressionRatio;

			Vector2 rr = new Vector2(v.x, v.z);
			var m = rr.magnitude;
			target = new Vector3(c * r * v.x / m, myY, c * r * v.z / m);


		}
		else
		{
			target = new Vector3(0, myY, 0);
		}
		v += (target - v) * speed;
		this.transform.position = v;
		updateLine();
	}

	void updateScale()
	{
		this.transform.localScale = new Vector3(
			(float)compressionRatio * scaleFactor,
			(float)compressionRatio * scaleFactor,
			(float)compressionRatio * scaleFactor
			);
	}

	void updateDisplay()
	{
		_textMesh.text =
			"id: " + modelHashID
			+ "\n message:\n " + commitMessage
			+ "\n c: " + compressionRatio;
	}

	public void Initialize()
	{
		updateDisplay();
		updateScale();

		_renderer = this.GetComponent<Renderer>();
		updateLine();
	}
}